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

	log.Printf("> Login server listening on: %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.Handle(conn)
	}
}

// Handle handles incoming connections and parses packets.
func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()

	// Log new connection and remote IP address.
	log.Printf("> New connection from IP %s", conn.RemoteAddr().String())

	// Read packet from connection.
	packet := make([]byte, 1024)
	length, err := conn.Read(packet)
	if err != nil {
		log.Printf("Error reading packet: %v", err)
		return
	}
	packet = packet[:length]

	// Create new packet instance and parse packet.
	p := network.Packet{
		Conn: conn,
	}
	if err := p.ParsePacket(packet); err != nil {
		log.Printf("Error parsing packet: %v", err)
		return
	}
}
