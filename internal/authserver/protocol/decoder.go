package protocol

import (
	"fmt"
	"io"
)

// Read opcode and choose appropriate Packet decoder

func DecodeRequest(reader io.Reader) (Request, error) {
	opcode := make([]byte, 1)
	_, err := io.ReadFull(reader, opcode)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode %v", err)
	}
	fmt.Println("opcode:", opcode)
	switch opcode[0] {
	case byte(CmdAuthLogonChallenge):
		return DecodeLogonChallengeRequest(reader)
	case byte(CmdAuthLogonProof):
		return DecodeLogonProofRequest(reader)
	case byte(CmdRealmList):
		return DecodeRealmListRequest(reader)
	default:
		fmt.Println("received unknown opcode:", opcode[0])
		return nil, fmt.Errorf("unknown opcode")
	}
}
