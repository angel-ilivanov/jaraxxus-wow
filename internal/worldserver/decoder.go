package worldserver

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

const (
	clientPacketSizeBytesCount = 2
	clientOpcodeBytesCount     = 4
)

var ErrUnknownOpcode = errors.New("received unknown opcode")

func DecodeRequest(reader io.Reader) (protocol.ClientMessage, error) {
	sizeBuffer := make([]byte, clientPacketSizeBytesCount)
	_, err := io.ReadFull(reader, sizeBuffer)
	if err != nil {
		return nil, fmt.Errorf("error reading size from client packet: %w", err)
	}
	size := binary.LittleEndian.Uint16(sizeBuffer)
	opcodeBuffer := make([]byte, clientOpcodeBytesCount)
	_, err = io.ReadFull(reader, opcodeBuffer)
	if err != nil {
		return nil, fmt.Errorf("error reading opcode from client packet: %w", err)
	}
	opcode := binary.LittleEndian.Uint32(opcodeBuffer)
	switch opcode {
	case uint32(protocol.OpcodeAuthSession):
		return protocol.DecodeAuthSession(size, reader)
	default:
		return nil, ErrUnknownOpcode
	}
}
