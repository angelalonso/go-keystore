package main

import (
	"strings"
  "net/http"
  "fmt"
  "io/ioutil"
  "path/filepath"
)

// Routes

// API route: get the key for a given user
func (a *App) GetKey(w http.ResponseWriter, r *http.Request) {
	user, _ := r.URL.Query()["user"]
	respondWithJSON(w, http.StatusOK, ReadStoredKey(user[0]))
}

// API route: get the key for all users
func (a *App) GetKeyList(w http.ResponseWriter, r *http.Request) {
	files, dir_err := ioutil.ReadDir(keysPath)
	if dir_err != nil {
		fmt.Printf("Reading keys directory failed with error %s\n", dir_err)
	} else {
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
}

// Other functions

// Reads an internally stored key
func ReadStoredKey(user string) string {
	dat, err := ioutil.ReadFile(keysPath + "/" + user + ".pub.pem")
	if err != nil {
		fmt.Printf("Reading from File failed with error %s\n", err)
	}
	return string(dat)
}
