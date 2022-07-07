package main

import (
	"go-test/config"
	"go-test/crypto"
	"go-test/server"
)

func main() {
	config.ConfigInstance.Load()
	crypto.RsaInstance.Load()
	server.ServerInstance.Listen()
}
