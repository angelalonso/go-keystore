package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// test getting a key
func TestGetKey(t *testing.T) {
	// do we get a key?
	// given a testing pair, can it decrypt an encrypted message?
	getkey_req, _ := http.NewRequest("GET", "/getkey/?user=test", nil)
	getkey_rsp := executeRequest(getkey_req)
	/*
		buf, _ := ioutil.ReadFile(keysPath + "/test.pub")
		getkey_expected := `"` + string(NormalizeNewline([]byte(buf))) + `"`

		NormalizeNewline([]byte(getkey_rsp.Body.String()))
		checkResponseCode(t, http.StatusOK, getkey_rsp.Code)
	*/
	pemString, _ := ioutil.ReadFile("./test")
	encString, _ := ioutil.ReadFile("./encrypted.txt")
	test, _ := ParseRsaPrivateKeyFromPemStr(string(pemString))

	//getkey_expected := ""
	var label = flag.String("label", "", "Label to use (filename by default)")
	getkey_expected, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, test, encString, []byte(*label))
	fmt.Println(encString)
	fmt.Println(err)

	checkResponseCode(t, http.StatusOK, getkey_rsp.Code)
	if body := getkey_rsp.Body.String(); body != string(getkey_expected) {
		t.Errorf("Expected "+string(getkey_expected)+". Got %s", body)
	} else {
		fmt.Println("- Test OK: Get an Existing Key")
	}
	// given a testing pair, can it decrypt an encrypted message?
}

// test key registry

func TestKeyStore(t *testing.T) {
	// list key files vs keys on memory
	// remove a key
	// clean up all keys
}

// https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
