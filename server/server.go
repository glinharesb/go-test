package server

import (
	"encoding/hex"
	"fmt"
	"go-test/network"
	"log"
	"net"
)

type Server struct{}

var ServerInstance Server = Server{}

func (s *Server) Listen() {
	listener, err := net.Listen("tcp", ":4597")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("[!] Server listening...")

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

	fmt.Printf("[!] Packet:\n%s", hex.Dump(packet))

	packetInstance := network.Packet{
		Conn: conn,
	}
	packetInstance.ParsePacket(packet)
}
