package main

import (
	"log"
	"net/http"

	"git.alterway.fr/multi-iaas-billing-exporter/app"
)

func main() {
	app := app.New()

	http.HandleFunc("/", app.Router.ServeHTTP)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Println(err)
	}
}
