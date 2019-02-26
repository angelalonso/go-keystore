package main

import (
	"strings"
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
)

// Routes

// API route: get the key for a given user
func (a *App) GetKey(w http.ResponseWriter, r *http.Request) {
	user, _ := r.URL.Query()["user"]
	respondWithJSON(w, http.StatusOK, storedKey2String(user[0]))
}

// API route: get the key for all users
func (a *App) GetKeyList(w http.ResponseWriter, r *http.Request) {
	files, dir_err := ioutil.ReadDir(keysPath)
	if dir_err != nil {
    res := "Reading keys directory failed with error "+ dir_err.Error()
    respondWithJSON(w, http.StatusInternalServerError, res)
	} else {
    res := []string{}
    for _, f := range files {
      if !f.IsDir() && strings.HasSuffix(f.Name(), ".pub.pem") {
        var extension = ".pub.pem"
        var name = f.Name()[0 : len(f.Name())-len(extension)]
        res = append(res, name)
      }
    }
    respondWithJSON(w, http.StatusOK, res)
  }
}

// Other functions

// Reads an internally stored key into a string
func storedKey2String(user string) string {
	dat, err := ioutil.ReadFile(keysPath + "/" + user + ".pub.pem")
	if err != nil {
		fmt.Printf("Reading from File failed with error %s\n", err)
	}
	return string(dat)
}

// Loads a public key from a file into a public key object
func keyFile2KeyObject(fileName string) *rsa.PublicKey {
  key_content, key_err := file2String(fileName)
	if key_err != nil {
		fmt.Println("Fatal error ", key_err.Error())
		fmt.Println(" - You can generate the key pair like this:")
		fmt.Println("     openssl genrsa -out testpriv.pem 4096")
		fmt.Println("     ssh-keygen -f testpriv.pem -e -m pem > testpub.pem")
		os.Exit(1)
	}
	size := len(key_content)
	pembytes := make([]byte, size)
	buffer := strings.NewReader(key_content)
	_, bufread_err := buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	publicKeyFileImported, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if bufread_err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	return publicKeyFileImported
}

// Returns contents of a file as string
func file2String(fileName string) (string, error) {
	dat, err := ioutil.ReadFile(fileName)
	return string(dat), err
}

// Loads a public key from a string
func pubKeyString2PubKeyObject(key string) *rsa.PublicKey {

	size := len(key)
	pembytes := make([]byte, size)
	buffer := strings.NewReader(key)
	_, err := buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	publicKeyImported, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	return publicKeyImported
}
