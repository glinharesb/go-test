package server

import (
	"fmt"
	"go-test/config"
	"go-test/network"
	"log"
	"net"
)

func Listen() {
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

		go Handle(conn)
	}
}

// handle handles incoming connections and parses packets
func Handle(conn net.Conn) {
	defer conn.Close()

	fmt.Println("")

	// log new connection and remote IP address
	log.Printf("> New connection from IP %s", conn.RemoteAddr().String())

	// read packet from connection
	packet := make([]byte, 1024)
	length, err := conn.Read(packet)
	if err != nil {
		log.Printf("Error reading packet: %v", err)
		return
	}
	packet = packet[:length]

	// create new packet instance and parse packet
	p := network.Packet{
		Conn: conn,
	}
	if err := p.ParsePacket(packet); err != nil {
		log.Printf("Error parsing packet: %v", err)
		return
	}
}
