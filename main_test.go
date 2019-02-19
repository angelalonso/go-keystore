package main

import (
	//"bytes"
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
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

func TestHealth(t *testing.T) {
	// Start the service
	//jsonData := map[string]string{"user": "Nic", "pubkey": "Raboy"}
	//jsonValue, _ := json.Marshal(jsonData)
	//response, err := http.Post("http://0.0.0.0:8400/key", "application/json", bytes.NewBuffer(jsonValue))
	// Health returns something proper
	health_req, _ := http.NewRequest("GET", "/health", nil)
	health_rsp := executeRequest(health_req)

	checkResponseCode(t, http.StatusOK, health_rsp.Code)

	if body := health_rsp.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
	/*
			health_out, err := http.Get("http://0.0.0.0:8400/health")
			if err != nil {
				t.Errorf("The HTTP request failed with error %s\n", err)
			} else {
				health_data, _ := ioutil.ReadAll(health_out.Body)
				fmt.Println(string(health_data))
			}

		  nohealth_out, nohealth_err := http.Get("http://0.0.0.0:8400/nohealth")
			if nohealth_err != nil {
				t.Errorf("The HTTP request failed with error %s\n", nohealth_err)
			} else {
				nohealth_data, _ := ioutil.ReadAll(nohealth_out.Body)
				fmt.Println(string(nohealth_data))
			}

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.

			/*
				if len(d) != 52 {
					t.Errorf("Expected deck length of 52, but got %v", len(d))
				}
	*/

}

// test adding a new key

func TestAddKey(t *testing.T) {
	// does it upload a key from a file?
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
