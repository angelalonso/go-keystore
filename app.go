package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

const keysPath = "./keys"

// Initializes the API server
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Initializes the routes for the API server
func (a *App) initializeRoutes() {
	//  curl http://0.0.0.0:8400/health
	a.Router.HandleFunc("/health", a.Health).Methods("GET")
	// curl -X POST -F 'file=@user1_key.pub' http://0.0.0.0:8400/addkey/\?user\=user1
	a.Router.Path("/addkey/").Methods("POST").HandlerFunc(a.AddPubKeyFile)
	// curl http://0.0.0.0:8400/getkey/\?user\=user1
	a.Router.HandleFunc("/getkey/", a.GetKey).Methods("GET")
	// curl http://0.0.0.0:8400/list
	a.Router.HandleFunc("/list", a.GetKeyList).Methods("GET")
}

// Runs the server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Routes

// API route: Health route just returns ok
func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "ok")
}

// Other functions

// Prepares response as a JSON object
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
