package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type authSessionWire1 struct {
	Build         uint32
	LoginServerId uint32
}

type authSessionWire2 struct {
	LoginServerType       uint32
	ClientSeed            [4]byte
	RegionId              uint32
	BattlegroupId         uint32
	RealmId               uint32
	DosResponse           uint64
	ClientProof           [20]byte
	DecompressedAddonInfo uint32
}

type AuthSessionRequest struct {
	username    string
	clientSeed  []byte
	clientProof []byte
}

func (a AuthSessionRequest) isClientMessage() {}

func DecodeAuthSession(size uint16, reader io.Reader) (AuthSessionRequest, error) {
	packet, err := readAuthSessionPacket(size, reader)
	if err != nil {
		return AuthSessionRequest{}, err
	}
	fmt.Println(packet)
	request, err := decodeAuthSessionRequest(packet)
	if err != nil {
		return AuthSessionRequest{}, fmt.Errorf("error decoding authSessionRequest: %w", err)
	}
	return request, nil
}

func readAuthSessionPacket(size uint16, reader io.Reader) ([]byte, error) {
	body := make([]byte, size-ClientOpcodeLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return nil, fmt.Errorf("error reading CMSG_AUTH_SESSION packet bytes: %w", err)
	}
	return body, nil
}

func decodeAuthSessionRequest(packetBytes []byte) (AuthSessionRequest, error) {
	var wire1 authSessionWire1
	reader := bytes.NewReader(packetBytes)
	err := binary.Read(reader, binary.LittleEndian, &wire1)
	if err != nil {
		fmt.Println(err)
		return AuthSessionRequest{}, fmt.Errorf("error reading wire1: %w", err)
	}
	name, err := readUsername(reader)
	if err != nil {
		return AuthSessionRequest{}, fmt.Errorf("error reading username: %w", err)
	}
	var wire2 authSessionWire2
	err = binary.Read(reader, binary.LittleEndian, &wire2)
	if err != nil {
		return AuthSessionRequest{}, fmt.Errorf("error reading wire2: %w", err)
	}

	//read rest of bytes to complete packet

	return AuthSessionRequest{
		username:    name,
		clientSeed:  wire2.ClientSeed[:],
		clientProof: wire2.ClientProof[:],
	}, nil
}

func readUsername(reader io.ByteReader) (string, error) {
	username := make([]byte, 0, MaxUsernameLength+1)
	currentByte, err := reader.ReadByte()
	if err != nil {
		return "", fmt.Errorf("error reading byte: %w", err)
	}
	for currentByte != 0x00 {
		username = append(username, currentByte)
		currentByte, err = reader.ReadByte()
		if err != nil {
			return "", fmt.Errorf("error reading byte: %w", err)
		}
	}
	username = append(username, 0x00)
	return string(username), nil
}
