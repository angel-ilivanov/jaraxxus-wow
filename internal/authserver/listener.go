package authserver

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func Start(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error establishing listener %s", err)
	}
	defer listener.Close()
	fmt.Println("Listening for TCP Connections")
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

	for {
		header := make([]byte, 4)
		_, err := io.ReadFull(conn, header) // Keeps reading until array is full
		if err != nil {
			log.Fatalf("error reading header %v", err)
			return
		}
		fmt.Println(header)
		size := binary.LittleEndian.Uint16(header[2:4]) // length without header
		body := make([]byte, size)
		_, err = io.ReadFull(conn, body)
		if err != nil {
			log.Fatalf("error reading body %v", err)
			return
		}
		fullPacket := append(header, body...)
		Parse(fullPacket)
	}
}
