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
}

func (a *App) Run(addr string) { 
  a.Router.HandleFunc("/health", Health).Methods("GET")
  a.Router.HandleFunc("/key", AddKey).Methods("POST")
  a.Router.Path("/keyfile").Methods("POST").HandlerFunc(AddKeyFile)
  a.Router.HandleFunc("/key/", GetKey).Methods("GET")
  log.Fatal(http.ListenAndServe(addr, a.Router))
}
