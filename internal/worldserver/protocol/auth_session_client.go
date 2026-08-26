package protocol

import (
	"fmt"
	"io"
)

type authSessionWire struct {
	build uint32
}

type authSessionRequest struct {
	username    string
	clientSeed  []byte
	clientProof []byte
}

func (a authSessionRequest) isClientMessage() {}

func DecodeAuthSession(size uint16, reader io.Reader) (ClientMessage, error) {
	packet, err := readAuthSessionPacket(size, reader)
	if err != nil {
		return nil, err
	}
	fmt.Println(packet)
	return nil, nil
}

func readAuthSessionPacket(size uint16, reader io.Reader) ([]byte, error) {
	body := make([]byte, size-4) //subtract opcode length
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return nil, fmt.Errorf("error reading CMSG_AUTH_SESSION packet bytes: %w", err)
	}
	fmt.Println(body)
	return body, nil
}
