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
	Username    string
	ClientSeed  []byte
	ClientProof []byte
}

func (a AuthSessionRequest) isClientMessage() {}

func (a AuthSessionRequest) Opcode() ClientOpcode {
	return OpcodeAuthSession
}

func DecodeAuthSession(size uint16, reader io.Reader) (AuthSessionRequest, error) {
	packetBody, err := readAuthSessionPacket(size, reader)
	if err != nil {
		return AuthSessionRequest{}, err
	}
	request, err := decodeAuthSessionRequest(packetBody)
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
	// Read static length part before username
	var wire1 authSessionWire1
	reader := bytes.NewReader(packetBytes)
	err := binary.Read(reader, binary.LittleEndian, &wire1)
	if err != nil {
		fmt.Println(err)
		return AuthSessionRequest{}, fmt.Errorf("error reading wire1: %w", err)
	}

	username, err := readUsername(reader)
	if err != nil {
		return AuthSessionRequest{}, fmt.Errorf("error reading username: %w", err)
	}

	// Read static length part after username
	var wire2 authSessionWire2
	err = binary.Read(reader, binary.LittleEndian, &wire2)
	if err != nil {
		return AuthSessionRequest{}, fmt.Errorf("error reading wire2: %w", err)
	}

	// Read rest of bytes to complete packet
	compressedAddonInfoLength := len(packetBytes) - 60 - (len(username) + 1)
	compressedAddonInfo := make([]byte, compressedAddonInfoLength)
	_, err = io.ReadFull(reader, compressedAddonInfo)
	if err != nil {
		return AuthSessionRequest{}, fmt.Errorf("error reading compressedAddonInfo: %w", err)
	}

	return AuthSessionRequest{
		Username:    username,
		ClientSeed:  wire2.ClientSeed[:],
		ClientProof: wire2.ClientProof[:],
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
	return string(username), nil
}
