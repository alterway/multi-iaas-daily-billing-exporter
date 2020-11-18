package app

import "github.com/gorilla/mux"

type App struct {
	Router *mux.Router
}

func New() *App {
	a := &App{
		Router: mux.NewRouter(),
	}

	a.initRoutes()
	return a
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/metrics", a.BillingHandler()).Methods("GET")
}
