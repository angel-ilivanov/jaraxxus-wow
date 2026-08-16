package authserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

func (s *Server) Start(ctx context.Context, port string) {
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
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection received")
	session := &AuthSession{}

	for {
		request, err := protocol.DecodeRequest(conn)
		if err != nil {
			return
		}
		packet, err := s.requestHandler.HandleRequest(ctx, session, request)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("packet length:", len(packet))
		bytesWritten, err := conn.Write(packet)
		fmt.Println("bytes written:", bytesWritten)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
