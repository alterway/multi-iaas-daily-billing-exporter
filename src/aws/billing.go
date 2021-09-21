package aws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

func CostExtract(w http.ResponseWriter) {
	//Load the Shared AWS Configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-1"),
		config.WithSharedCredentialsFiles([]string{"/var/secrets/aws/credentials"}),
		config.WithSharedConfigFiles([]string{"/var/secrets/aws/config"}),
	)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	// Create a Cost Explore service client
	client := costexplorer.NewFromConfig(cfg)

	// Input paramters
	metrics := []string{"AmortizedCost"}
	granularity := "DAILY"

	today := time.Now()
	yesterday := time.Now().Add(-48 * time.Hour)
	yYear, yMonth, yDay := yesterday.UTC().Date()
	tYear, tMonth, tDay := today.UTC().Date()
	startDate := fmt.Sprintf("%d-%02d-%02d", yYear, yMonth, yDay)
	endDate := fmt.Sprintf("%d-%02d-%02d", tYear, tMonth, tDay)

	// Input for request
	costInput := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startDate),
			End:   aws.String(endDate),
		},
		Granularity: types.Granularity(granularity),
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
		Metrics: metrics,
	}

	//Get cost
	result, err := client.GetCostAndUsage(context.TODO(), costInput)

	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	resultByTime := result.ResultsByTime
	for _, resultByDay := range resultByTime {
		groupsResult := resultByDay.Groups
		timePeriod := resultByDay.TimePeriod
		if len(groupsResult) > 0 {
			for _, groupResult := range groupsResult {
				fmt.Fprintln(w, prometheusExport(groupResult, *timePeriod))
			}
		}
	}
}

func prometheusExport(group types.Group, timePeriod types.DateInterval) string {
	serviceName := group.Keys
	cost := group.Metrics["AmortizedCost"]
	costAmount := aws.ToString(cost.Amount)
	costCurrency := aws.ToString(cost.Unit)
	startDate := aws.ToString(timePeriod.Start)

	return fmt.Sprintf(
		"cloud_cost{provider=\"aws\", currency=\"%s\", service=\"%s\", startDate=\"%s\"} %s",
		costCurrency,
		serviceName,
		startDate,
		costAmount)
}
