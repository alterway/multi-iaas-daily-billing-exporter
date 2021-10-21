package app

import (
	"log"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func New() *App {
	a := &App{
		Router: mux.NewRouter(),
	}

	log.Printf("Starting app...\n")
	a.initRoutes()
	a.initTargets()
	return a
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/metrics", a.BillingHandler()).Methods("GET")
}

func (a *App) initTargets() {
	a.TargetHandler()
}
