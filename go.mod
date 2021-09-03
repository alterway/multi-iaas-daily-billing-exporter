module git.alterway.fr/multi-iaas-billing-exporter

go 1.15

require (
	cloud.google.com/go/bigquery v1.21.0
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.3.3
	github.com/aws/aws-sdk-go-v2/config v1.1.6
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.3.1
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
	google.golang.org/api v0.54.0
)
