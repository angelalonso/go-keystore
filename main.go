package main

import ()

type Key struct {
	User   string `json:"user"`
	PubKey string `json:"pubkey"`
}

var a App

var keys []Key

func main() {
	a = App{}
	a.Initialize()
	a.Run(":8400")
}
