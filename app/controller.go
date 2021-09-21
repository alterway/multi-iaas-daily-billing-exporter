package app

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"git.alterway.fr/multi-iaas-billing-exporter/src/gcp"
)

func (a *App) BillingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gcp.ReportHandler(w)
		// aws.CostExtract(w)
	}
}

func (a *App) TargetHandler() {
	log.Printf("Configuring Prometheus scraping. This can take up to 2 minutes on top of scraping interval...\n")
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				resp, err := http.Get("http://prometheus-prometheus-oper-prometheus.monitoring.svc.cluster.local:9090/api/v1/targets")
				if err != nil {
					log.Fatalln(err)

				}
				defer resp.Body.Close()

				data := new(TargetResponse)
				json.NewDecoder(resp.Body).Decode(data)

				for _, target := range data.Data.ActiveTargets {

					if target.ScrapePool == "billing/multi-iaas-daily-billing-exporter-svc-monitor/0" && target.Health == "up" {
						log.Printf("Billing target detected\n")
						log.Printf("Target health: %v\n", target.Health)
						log.Printf("Prometheus will now start to scrape metrics\n")
						close(quit)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
