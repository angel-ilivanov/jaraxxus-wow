package protocol

import (
	"fmt"
	"io"
	"net"
)

// Read opcode and choose appropriate Packet decoder

func ReadClientMessage(conn net.Conn) (ClientMessage, error) {
	opcode := make([]byte, 1)
	_, err := io.ReadFull(conn, opcode)
	if err != nil {
		return LogonChallengeClientRequest{}, fmt.Errorf("error reading opcode %v", err)
	}
	fmt.Println("opcode: ")
	fmt.Println(opcode)
	switch opcode[0] {
	case 0x00:
		return ParseLogonChallengeClient(conn)
	default:
		return LogonChallengeClientRequest{}, fmt.Errorf("unknown opcode")
	}
}
