package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//	"os"
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
	router.
		Path("/keyfile").
		Methods("POST").
		HandlerFunc(AddKeyFile)
	router.HandleFunc("/key/", GetKey).Methods("GET")
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

func AddKeyFile(w http.ResponseWriter, r *http.Request) {
	//https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
	var Buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	contents := Buf.String()
	fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	return
}

func WriteKey(user string, key string) {
	err := ioutil.WriteFile(user+".pub", []byte(key), 0644)
	if err != nil {
		fmt.Printf("Writing to File failed with error %s\n", err)
	}
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	user, _ := r.URL.Query()["user"]
	fmt.Println(ReadKey(user[0]))
}

func ReadKey(user string) string {
	dat, err := ioutil.ReadFile(user + ".pub")
	if err != nil {
		fmt.Printf("Reading from File failed with error %s\n", err)
	}
	return string(dat)

}
