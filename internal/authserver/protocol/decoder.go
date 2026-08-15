package protocol

import (
	"fmt"
	"io"
)

// Read opcode and choose appropriate Packet decoder

func DecodeRequest(conn io.Reader) (ClientMessage, error) {
	opcode := make([]byte, 1)
	_, err := io.ReadFull(conn, opcode)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode %v", err)
	}
	fmt.Println("opcode:", opcode)
	switch opcode[0] {
	case 0x00:
		return DecodeLogonChallengeRequest(conn)
	case 0x01:
		fmt.Println("Received proof packet from client, handling not implemented")
		return nil, fmt.Errorf("logon proof handling not implemented")
	default:
		return nil, fmt.Errorf("unknown opcode")
	}
}
