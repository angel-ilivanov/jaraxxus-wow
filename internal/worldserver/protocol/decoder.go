package protocol

import (
	"errors"
	"fmt"
	"io"
)

var ErrUnknownOpcode = errors.New("received unknown opcode")

func DecodeRequest(header ClientHeader, reader io.Reader) (ClientMessage, error) {
	switch header.Opcode {
	case ClientOpcodeAuthSession:
		return DecodeAuthSession(header.PacketSize, reader)
	case ClientOpcodeAccountTimesReady:
		return AccountDataTimesRequest{}, nil
	case ClientOpcodeCharEnum:
		return CharEnumRequest{}, nil
	case ClientOpcodeCharCreate:
		return DecodeCharCreate(header.PacketSize, reader)
	default:
		return nil, fmt.Errorf("%w: 0x%02X", ErrUnknownOpcode, header.Opcode)
	}
}
