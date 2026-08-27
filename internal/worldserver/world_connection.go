package worldserver

import (
	"errors"
	"fmt"
	"net"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

const (
	largeHeaderFlag      = 0x80     // 0b10000000, bitwise or with MSB of size field to set flag
	largeHeaderThreshold = 0x7FFF   // 15 bits (16th bit is reserved for largeHeaderFlag)
	sizeFieldMaxValue    = 0x7FFFFF // 23 size bits + 1 flag bit
)

var ErrHeaderSizeTooLarge = errors.New("header size is too large")

type WorldConnection struct {
	headerCipher *srp6.HeaderCipher
	connection   net.Conn
}

func (c *WorldConnection) WriteMessage(message protocol.ServerMessage) error {
	body := message.EncodeBody()
	header, err := constructHeader(message.Opcode(), len(body))
	if err != nil {
		return fmt.Errorf("error constructing header: %w", err)
	}

	if c.headerCipher != nil { // If cipher is absent, send packet as plain text
		err = c.headerCipher.EncryptHeader(header)
		if err != nil {
			return fmt.Errorf("error encrypting header: %w", err)
		}
	}

	packet := append(header, body...)
	bytesWritten, err := c.connection.Write(packet)
	if err != nil {
		return fmt.Errorf("error writing packet: %w", err)
	}
	if bytesWritten != len(packet) {
		return fmt.Errorf("bytesWritten does not match packet length, expected %d but got %d", len(packet), bytesWritten)
	}
	return nil
}

func (c *WorldConnection) ReadMessage() (protocol.ClientMessage, error) {
	// 1. decode header if encryption is enabled
	// 2. send to decoder
	// TODO: IMPLEMENT THIS
	panic("implement me")
}

func constructHeader(opcode protocol.ServerOpcode, bodyLength int) ([]byte, error) {
	size := bodyLength + protocol.ServerOpcodeLength

	if size > sizeFieldMaxValue {
		return nil, ErrHeaderSizeTooLarge
	}
	var header []byte
	if size <= largeHeaderThreshold {
		// Normal header, 4 bytes
		// Size (big endian, 2 bytes), Opcode (little endian, 2 bytes)
		header = []byte{
			byte(size >> 8),
			byte(size),
			byte(opcode),
			byte(opcode >> 8),
		}
	} else {
		// Large header, 5 bytes
		// 3 byte size, big endian (set flag on MSB)
		// 2 byte opcode, little endian
		header = []byte{
			byte(size>>16) | largeHeaderFlag,
			byte(size >> 8),
			byte(size),
			byte(opcode),
			byte(opcode >> 8),
		}
	}

	return header, nil
}

func NewWorldConnection(conn net.Conn) *WorldConnection {
	return &WorldConnection{
		headerCipher: nil,
		connection:   conn,
	}
}

func (c *WorldConnection) EnableEncryption(sessionKey []byte) error {
	c.headerCipher = &srp6.HeaderCipher{}
	err := c.headerCipher.Init(sessionKey)
	if err != nil {
		return err
	}
	return nil
}
