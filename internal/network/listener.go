package network

import (
	"fmt"
	"net"
	"wow-server/internal/auth"
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

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("Error reading bytes from connection: %s", err)
		return
	}
	fmt.Printf("Full byte string: %x\n", buf[:n])
	auth.ParsePacket(buf[:n])
}
