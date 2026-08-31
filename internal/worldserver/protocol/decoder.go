package protocol

import (
	"errors"
	"fmt"
	"io"
)

var ErrUnknownOpcode = errors.New("received unknown opcode")

func DecodeRequest(header ClientHeader, reader io.Reader) (ClientMessage, error) {
	switch header.Opcode {
	case OpcodeAuthSession:
		return DecodeAuthSession(header.PacketSize, reader)
	case OpcodeAccountTimesReady:
		return AccountDataTimesRequest{}, nil
	case OpcodeCharEnum:
		return nil, fmt.Errorf("OpcodeCharEnum not yet implemented")
	default:
		return nil, fmt.Errorf("%w: 0x%02X", ErrUnknownOpcode, header.Opcode)
	}
}
