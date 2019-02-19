package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/health", a.Health).Methods("GET")
	a.Router.HandleFunc("/key", AddKey).Methods("POST")
	a.Router.Path("/keyfile").Methods("POST").HandlerFunc(AddKeyFile)
	a.Router.HandleFunc("/key/", GetKey).Methods("GET")
}

func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "ok")
}
