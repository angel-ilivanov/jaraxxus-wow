package protocol

import (
	"fmt"
	"io"
)

// Read opcode and choose appropriate Packet decoder

func ReadClientMessage(conn io.Reader) (ClientMessage, error) {
	opcode := make([]byte, 1)
	_, err := io.ReadFull(conn, opcode)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode %v", err)
	}
	fmt.Println("opcode: ")
	fmt.Println(opcode)
	switch opcode[0] {
	case 0x00:
		return ParseLogonChallengeClient(conn)
	default:
		return nil, fmt.Errorf("unknown opcode")
	}
}
