package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Key struct {
	User   string `json:"user"`
	PubKey string `json:"pubkey"`
}

var keys []Key

func main() {
	router := mux.NewRouter()
	// call with:
	//  curl -H "Content-Type: application/json" -d '{"user":"hey","pubkey":"hoooo"}' -X POST http://127.0.0.1:8400/key
	router.HandleFunc("/health", Health).Methods("GET")
	router.HandleFunc("/key", AddKey).Methods("POST")
	router.HandleFunc("/key/{userkey}", GetKey).Methods("GET")
	log.Fatal(http.ListenAndServe(":8400", router))
}

func Health(w http.ResponseWriter, r *http.Request) {
	return
}

func AddKey(w http.ResponseWriter, r *http.Request) {
	var newkey Key
	_ = json.NewDecoder(r.Body).Decode(&newkey)
	// TODO: check error
	WriteKey(newkey.User, newkey.PubKey)

}

func WriteKey(user string, key string) {
	err := ioutil.WriteFile(user+".pub", []byte(key), 0644)
	if err != nil {
		fmt.Printf("Writing to File failed with error %s\n", err)
	}

}

func GetKey(w http.ResponseWriter, r *http.Request) {

}
