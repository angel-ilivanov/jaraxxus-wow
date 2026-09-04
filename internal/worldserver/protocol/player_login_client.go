package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type PlayerLoginRequest struct {
	charGUID uint64
}

func (p PlayerLoginRequest) isClientMessage() {}

func (p PlayerLoginRequest) Opcode() ClientOpcode {
	return ClientOpcodePlayerLogin
}

func DecodePlayerLogin(reader io.Reader) (PlayerLoginRequest, error) {
	var request PlayerLoginRequest
	err := binary.Read(reader, binary.LittleEndian, &request)
	if err != nil {
		return PlayerLoginRequest{}, fmt.Errorf("error reading character GUID: %w", err)
	}
	return request, nil
}
