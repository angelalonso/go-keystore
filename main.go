package main

import ()

type Key struct {
	User   string `json:"user"`
	PubKey string `json:"pubkey"`
}
// Variables used all across the application
var a App

var keys []Key

// main is just the initialization of the API server
func main() {
	a = App{}
	a.Initialize()
	a.Run(":8400")
}
