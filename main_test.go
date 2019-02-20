package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
	file, err := os.Open("./test.pub")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	addkey_req, _ := http.NewRequest("POST", "/addkey/?user=test", file)
	addkey_rsp := executeRequest(addkey_req)

	checkResponseCode(t, http.StatusOK, addkey_rsp.Code)

	if body := addkey_rsp.Body.String(); body != `"ok"` {
		t.Errorf("Expected an 'ok' message. Got %s", body)
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

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
