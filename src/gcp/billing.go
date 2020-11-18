package gcp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

type Cost struct {
	Amount   float64
	Currency string
}

type Report struct {
	Cost       Cost
	Service    string
	ProjectID  string
	DateReport string
}

type Row struct {
	Total_cost  float64
	Date_report string
	Description string
	Id          string
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getEnv(name string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return goDotEnvVariable(name)
}

func (r *Report) prometheus() string {
	return fmt.Sprintf(
		"gcp_cost{currency=\"%s\", service=\"%s\", dateReport=\"%s\", project=\"%s\"} %.2f",
		r.Cost.Currency,
		r.Service,
		r.DateReport,
		r.ProjectID,
		r.Cost.Amount)
}

func ReportHandler(w http.ResponseWriter) {
	ctx := context.Background()
	projectID := getEnv("GCP_PROJECTID")
	queryTable := getEnv("GCP_TABLE")
	service, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer service.Close()

	q_daily := service.Query(
		"SELECT STRING(DATE((usage_start_time))) AS date_report, service.description, project.id, SUM(cost) AS total_cost FROM  " + queryTable +
			"WHERE DATE(usage_start_time) >= DATE_ADD(CURRENT_DATE(), INTERVAL -1 DAY) " +
			"GROUP BY date_report, service.description, project.id " +
			"ORDER BY date_report")

	// Run the query and print results when the query job is completed.
	job, err := q_daily.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := status.Err(); err != nil {
		log.Fatal(err)
	}
	it, err := job.Read(ctx)
	for {
		var row Row
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		report := Report{
			Cost: Cost{
				Currency: "EUR",
				Amount:   row.Total_cost,
			},
			Service:    row.Description,
			ProjectID:  row.Id,
			DateReport: row.Date_report,
		}
		fmt.Fprintln(w, report.prometheus())
	}
}
