package protocol

import (
	"fmt"
	"io"
)

// Read opcode and choose appropriate Packet decoder

func DecodeRequest(conn io.Reader) (Request, error) {
	opcode := make([]byte, 1)
	_, err := io.ReadFull(conn, opcode)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode %v", err)
	}
	fmt.Println("opcode:", opcode)
	switch opcode[0] {
	case byte(CmdAuthLogonChallenge):
		return DecodeLogonChallengeRequest(conn)
	case byte(CmdAuthLogonProof):
		return DecodeLogonProofRequest(conn)
	case 0x10:
		fmt.Println("received Realm List request, handling not yet implemented")
		return nil, fmt.Errorf("handling for realm list request not yet implemented")
	default:
		fmt.Println("received unknown opcode:", opcode[0])
		return nil, fmt.Errorf("unknown opcode")
	}
}
