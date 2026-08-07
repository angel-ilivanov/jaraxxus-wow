package main

import (
	"fmt"
	"net"
)

var port = ":3724"

func main() {
	//tcp server listening on port 3724
	//client connects -> spin up a goroutine to handle session

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting TCP server: ", err)
	}
	defer listener.Close()
	fmt.Println("Listening for TCP connections")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error")
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection received!")
}
