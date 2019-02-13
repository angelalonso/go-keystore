package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewKey(t *testing.T) {
	// Start the service
	jsonData := map[string]string{"user": "Nic", "pubkey": "Raboy"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://0.0.0.0:8400/key", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.

	// does it upload a key from a file?
	// does the key have a username, does it alert if not?
	// does the key have a public key, does it alert if not?
	// does it save it to a file?
	/*
		if len(d) != 52 {
			t.Errorf("Expected deck length of 52, but got %v", len(d))
		}
	*/

}

func TestGetKey(t *testing.T) {
	// given a testing pair, can it decrypt an encrypted message?
	/*
		filename := "_decktesting"
		os.Remove(filename)

		deck := newDeck()
		deck.saveToFile(filename)

		loadedDeck := newDeckFromFile(filename)
		if len(loadedDeck) != 52 {
			t.Errorf("Expected deck length of 52, but got %v", len(loadedDeck))
		}

		os.Remove(filename)
	*/
}
