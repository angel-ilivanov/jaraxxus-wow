package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	clientPacketSizeBytesCount = 2
	clientOpcodeBytesCount     = 4
)

var ErrUnknownOpcode = errors.New("received unknown opcode")

func DecodeRequest(reader io.Reader) (ClientMessage, error) {
	sizeBuffer := make([]byte, clientPacketSizeBytesCount)
	_, err := io.ReadFull(reader, sizeBuffer)
	if err != nil {
		return nil, fmt.Errorf("error reading size from client packet: %w", err)
	}
	size := binary.BigEndian.Uint16(sizeBuffer)

	opcodeBuffer := make([]byte, clientOpcodeBytesCount)
	_, err = io.ReadFull(reader, opcodeBuffer)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode from client packet: %w", err)
	}
	opcode := binary.LittleEndian.Uint32(opcodeBuffer)

	switch opcode {
	case uint32(OpcodeAuthSession):
		return DecodeAuthSession(size, reader)
	case uint32(OpcodeCharEnum):
		return nil, fmt.Errorf("OpcodeCharEnum not yet implemented")
	default:
		return nil, fmt.Errorf("%w: 0x%02X", ErrUnknownOpcode, opcode)
	}
}
