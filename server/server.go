package server

import (
	"fmt"
	"go-test/config"
	"go-test/network"
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

func Listen() error {
	address := fmt.Sprintf("%s:%d", config.GetConfig().LoginIp, config.GetConfig().LoginPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	logrus.Infof("login server listening on: %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	defer conn.Close()

	packet := make([]byte, 1024)
	length, err := conn.Read(packet)
	if err != nil {
		log.Printf("error reading packet: %v", err)
		return
	}
	packet = packet[:length]

	p := network.Packet{
		Conn: conn,
	}
	if err := p.ParsePacket(packet); err != nil {
		log.Printf("error parsing packet: %v", err)
		return
	}
}
