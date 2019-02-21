package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// test health, nohealth

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)

}

// test the health URL
func TestHealth(t *testing.T) {
	// Health returns something proper
	health_req, _ := http.NewRequest("GET", "/health", nil)
	health_rsp := executeRequest(health_req)

	checkResponseCode(t, http.StatusOK, health_rsp.Code)

	if body := health_rsp.Body.String(); body != `"ok"` {
		t.Errorf("Expected an 'ok' message. Got %s", body)
	}
}

// test adding a new key
func TestAddKey(t *testing.T) {
	// Delete the keys folder, does it recreate it?
	os.RemoveAll(keysPath)
	// does it upload a key from a file?
	addkey_body, writer := MultipartUpload("./test.pub")
	addkey_req, _ := http.NewRequest("POST", "/addkey/?user=test", addkey_body)
	addkey_req.Header.Set("Content-Type", writer.FormDataContentType())

	addkey_rsp := executeRequest(addkey_req)

	checkResponseCode(t, http.StatusOK, addkey_rsp.Code)
	if body := addkey_rsp.Body.String(); body != `` {
		t.Errorf("Expected an empty message. Got %s", body)
	}
	// does the key have a username, does it alert if not?
	// does the key have a public key, does it alert if not?
	// does it save it to a file?
	// given a testing pair, can it decrypt an encrypted message?
}

// test getting a key
func TestGetKey(t *testing.T) {
	// do we get a key?
	// given a testing pair, can it decrypt an encrypted message?
}

// test key registry

func TestKeyStore(t *testing.T) {
	// list key files vs keys on memory
	// remove a key
	// clean up all keys
}

//https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func MultipartUpload(path string) (*bytes.Buffer, *multipart.Writer) {
	paramName := "file"
	params := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
    fmt.Println("Error opening the test key file. Have you created your keypair?")
    fmt.Println("  just run ssh-keygen -f test -t rsa -N ''")
		panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(part, file)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}
	return body, writer
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
