package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type App struct {
	Router *mux.Router
}

const keysPath = "./keys"

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

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

func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "ok")
}

func loadPublicPemKey(fileName string) *rsa.PublicKey {

	publicKeyFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		fmt.Println(" - You can generate the key pair like this:")
		fmt.Println("     openssl genrsa -out testpriv.pem 4096")
		fmt.Println("     ssh-keygen -f testpriv.pem -e -m pem > testpub.pem")
		os.Exit(1)
	}

	pemfileinfo, _ := publicKeyFile.Stat()

	size := pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	publicKeyFile.Close()
	publicKeyFileImported, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	return publicKeyFileImported
}

func (a *App) GetKey(w http.ResponseWriter, r *http.Request) {
	user, _ := r.URL.Query()["user"]
	fmt.Println(ReadKey(user[0]))
	respondWithJSON(w, http.StatusOK, ReadKey(user[0]))
}

func (a *App) GetKeyList(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir(keysPath)
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".pub.pem") {
			var extension = filepath.Ext(f.Name())
			var name = f.Name()[0 : len(f.Name())-len(extension)]
			res = append(res, name)
		}
	}
	fmt.Println(res)
}

func ReadKey(user string) string {
	dat, err := ioutil.ReadFile(keysPath + "/" + user + ".pub.pem")
	if err != nil {
		fmt.Printf("Reading from File failed with error %s\n", err)
	}
	return string(dat)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
