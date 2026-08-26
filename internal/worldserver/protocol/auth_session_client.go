package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type authSessionWire1 struct {
	Build         uint32
	LoginServerId uint32
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
	decodeAuthSessionWire(packet)
	return nil, nil
}

func readAuthSessionPacket(size uint16, reader io.Reader) ([]byte, error) {
	body := make([]byte, size-4) //subtract opcode length
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return nil, fmt.Errorf("error reading CMSG_AUTH_SESSION packet bytes: %w", err)
	}
	return body, nil
}

func decodeAuthSessionWire(packetBytes []byte) {
	var wire1 authSessionWire1
	reader := bytes.NewReader(packetBytes)
	err := binary.Read(reader, binary.LittleEndian, &wire1)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = readUsername(reader)
	if err != nil {
		return
	}
}

func readUsername(reader io.Reader) (string, error) {
	br := bufio.NewReader(reader)
	username, err := br.ReadBytes(0x00)
	if err != nil {
		return "", fmt.Errorf("error reading username: %w", err)
	}
	return string(username), nil
}
