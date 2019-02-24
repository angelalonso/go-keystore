package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//TODO: OBSOLETE?
func (a *App) AddKeyFile(w http.ResponseWriter, r *http.Request) {
	//https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
	user, _ := r.URL.Query()["user"]
	var Buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	strings.Split(header.Filename, ".")
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	contents := Buf.String()
	if len(user) > 0 {
		WriteKey(user[0], contents)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		Buf.Reset()
		// do something else
		// etc write header
		respondWithJSON(w, http.StatusOK, "Public key for user "+user[0]+" saved")
	} else {
		respondWithJSON(w, http.StatusBadRequest, "Error: no user was provided. Try ?user=username")
	}

	//return
}

func (a *App) AddPubKeyFile(w http.ResponseWriter, r *http.Request) {
	//https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
	user, _ := r.URL.Query()["user"]
	var Buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	strings.Split(header.Filename, ".")
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	if len(user) > 0 {
		// delimiter found at https://stackoverflow.com/questions/47190864/bufio-readbytes-buffer-size
		bufbytes, _ := Buf.ReadBytes('\x04')
		data, _ := pem.Decode(bufbytes)
		publicKeyFileReceived, err := x509.ParsePKCS1PublicKey(data.Bytes)
		if err != nil {
			fmt.Println("Fatal error on Adding Public Key ", err.Error())
			os.Exit(1)
		}
		WritePubKey(user[0], publicKeyFileReceived)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		Buf.Reset()
		// do something else
		// etc write header
		respondWithJSON(w, http.StatusOK, "Public key for user "+user[0]+" saved")
	} else {
		respondWithJSON(w, http.StatusBadRequest, "Error: no user was provided. Try ?user=username")
	}
}

func WriteKey(user string, key string) {
	if _, err_pathexists := os.Stat(keysPath); os.IsNotExist(err_pathexists) {
		os.MkdirAll(keysPath, os.ModePerm)
	}
	err := ioutil.WriteFile(keysPath+"/"+user+".pub", []byte(key), 0644)
	if err != nil {
		fmt.Printf("Writing to File failed with error %s\n", err)
	}
}

func WritePubKey(user string, pubkey *rsa.PublicKey) {
	//converts an RSA public key to PKCS#1, ASN.1 DER form.
	if _, err_pathexists := os.Stat(keysPath); os.IsNotExist(err_pathexists) {
		os.MkdirAll(keysPath, os.ModePerm)
	}
	var pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubkey),
	}
	fileName := keysPath + "/" + user + ".pub.pem"
	pemfile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	defer pemfile.Close()
	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
