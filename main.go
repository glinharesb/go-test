package main

import (
	"go-test/crypto"
	"go-test/server"
)

func main() {
	crypto.RsaInstance.Load()
	server.ServerInstance.Listen()
}
