package main

import (
	"go-test/config"
	"go-test/crypto"
	"go-test/database"
	"go-test/server"
)

func main() {
	config.ConfigInstance.Load()

	if database.DatabaseInstance.Load() {
		// account := database.DatabaseInstance.LoadAccountByName("4050344")
		// fmt.Println(account)
		// os.Exit(0)
		defer database.DatabaseInstance.Connection.Close()
	}

	crypto.RsaInstance.Load()
	server.ServerInstance.Listen()
}
