package authserver

import (
	"fmt"
	"log"
	"net"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

func (s *Server) Start(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error establishing listener %v", err)
	}
	defer listener.Close()
	fmt.Println("Listening for TCP Connections")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection received")

	for {
		_, err := protocol.ReadClientMessage(conn) //request, err
		if err != nil {
			return
		}
		//response = server.HandleMessage(session, request)
	}
}
