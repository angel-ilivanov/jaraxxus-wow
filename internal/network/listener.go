package network

import (
	"fmt"
	"net"
)

func Start(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error establishing listener %s", err)
	}
	defer listener.Close()
	fmt.Println("Listening for TCP Connetcions")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection received")
}
