package protocol

import (
	"errors"
	"fmt"
	"io"
)

// Read opcode and choose appropriate Packet decoder

var ErrUnknownOpcode = errors.New("unknown opcode")

func DecodeRequest(reader io.Reader) (Request, error) {
	opcode := make([]byte, 1)
	_, err := io.ReadFull(reader, opcode)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode %w", err)
	}
	switch opcode[0] {
	case byte(CmdAuthLogonChallenge):
		return DecodeLogonChallengeRequest(reader)
	case byte(CmdAuthLogonProof):
		return DecodeLogonProofRequest(reader)
	case byte(CmdRealmList):
		return DecodeRealmListRequest(reader)
	default:
		return nil, fmt.Errorf("%w: 0x%02X", ErrUnknownOpcode, opcode[0])
	}
}
