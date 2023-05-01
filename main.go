package main

import (
	"go-test/config"
	"go-test/crypto"
	"go-test/database"
	"go-test/server"
)

func main() {
	var err error

	config.Load()

	err = database.Connect()
	if err != nil {
		panic(err)
	}

	crypto.LoadRsa()

	server.Listen()
}
