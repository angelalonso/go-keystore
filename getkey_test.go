package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
  "io/ioutil"
)

// test getting a key

func TestGetKey(t *testing.T) {
	// do we get a key?

	manualKey := keyFile2KeyObject("./testpub.pem")

	getkey_req, _ := http.NewRequest("GET", "/getkey/?user=test", nil)
	getkey_rsp := executeRequest(getkey_req)

	checkResponseCode(t, http.StatusOK, getkey_rsp.Code)

	body := strings.Replace(getkey_rsp.Body.String(), `\n`, "\n", -1)
	body = body[1 : len(body)-1]
	bodyKey := pubKeyString2PubKeyObject(body)

	if bodyKey.N.Cmp(manualKey.N) != 0 || bodyKey.E != manualKey.E {
		// https://stackoverflow.com/questions/32042989/go-lang-differentiate-n-and-line-break?rq=1
		t.Errorf("Expected:\n %s and %s\n Got:\n %s and %s", manualKey.N.String(), string(manualKey.E), bodyKey.N.String(), string(bodyKey.E))
	} else {
		fmt.Println("- Test OK: upload key without keys folder")
	}

	// Let's test decryption
	// given a testing pair, can it decrypt an encrypted message?
}

// test key registry

func TestKeyStore(t *testing.T) {
	// list key files vs keys on memory
	list_req, _ := http.NewRequest("GET", "/list", nil)
	list_rsp := executeRequest(list_req)

	checkResponseCode(t, http.StatusOK, list_rsp.Code)

	dir_read, _ := ioutil.ReadDir(keysPath)
  expected_files := []string{}
  for _, f := range dir_read {
    if !f.IsDir() && strings.HasSuffix(f.Name(), ".pub.pem") {
      var extension = ".pub.pem"
      var name = f.Name()[0 : len(f.Name())-len(extension)]
      // very hacky but well
      expected_files = append(expected_files, string(`"`+name+`"`))
    }
  }
  // very hacky too
  expected_files_string := `[`+strings.Join(expected_files," ")+`]`
  if list_rsp.Body.String() != expected_files_string {
    t.Errorf("Expected:\n %s\n Got:\n %s", expected_files_string, list_rsp.Body.String())
	} else {
		fmt.Println("- Test OK: get all keys stored in the keys folder")
  }
	// remove a key
	// clean up all keys
}

func loadRSAPrivatePemKey(fileName string) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		fmt.Println(" - You can generate the key pair like this:")
		fmt.Println("     openssl genrsa -out testpriv.pem 4096")
		fmt.Println("     ssh-keygen -f testpriv.pem -e -m pem > testpub.pem")
		os.Exit(1)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	data, _ := pem.Decode([]byte(pembytes))
	privateKeyFile.Close()
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	return privateKeyImported
}
