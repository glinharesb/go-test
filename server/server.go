package server

import (
	"fmt"
	"go-test/config"
	"go-test/network"
	"log"
	"net"
)

type Server struct{}

var ServerInstance Server = Server{}

func (s *Server) Listen() {
	address := fmt.Sprintf("%s:%d", config.ConfigInstance.LoginIp, config.ConfigInstance.LoginPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("[!] Login server listening on: %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()

	packet := make([]byte, 1024)

	length, err := conn.Read(packet)
	if err != nil {
		log.Fatal(err)
	}
	packet = packet[:length]

	packetInstance := network.Packet{
		Conn: conn,
	}
	packetInstance.ParsePacket(packet)
}
