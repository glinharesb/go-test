package main

import (
	"go-test/config"
	"go-test/crypto"
	"go-test/database"
	"go-test/server"
	"log"
	"os"
)

func main() {
	// Create or open the log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set the log output to the file
	// log.SetOutput(file)

	config.ConfigInstance.Load()

	database.DatabaseInstance.Connect()
	defer database.DatabaseInstance.Connection.Close()

	crypto.RsaInstance.Load()
	server.ServerInstance.Listen()
}
