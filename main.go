package main

import (
	"go-test/config"
	"go-test/crypto"
	"go-test/database"
	"go-test/server"

	"github.com/sirupsen/logrus"
)

func main() {
	// Set up logrus
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// logrus.SetOutput(os.Stdout)

	// Load application settings
	if err := config.Load(); err != nil {
		logrus.WithError(err).Fatal("failed to load application settings")
	}

	// Connect to the database
	if err := database.Connect(); err != nil {
		logrus.WithError(err).Fatal("failed to connect to the database")
	}

	// Load RSA private key
	if err := crypto.LoadRsa(); err != nil {
		logrus.WithError(err).Fatal("failed to load RSA private key")
	}

	// Start the server
	if err := server.Listen(); err != nil {
		logrus.WithError(err).Fatal("failed to start server")
	}
}
