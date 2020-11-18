package app

import (
	"net/http"

	"git.alterway.fr/multi-iaas-billing-exporter/src/aws"
	"git.alterway.fr/multi-iaas-billing-exporter/src/gcp"
)

func (a *App) BillingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gcp.ReportHandler(w)
		aws.CostExtract(w)
	}
}
