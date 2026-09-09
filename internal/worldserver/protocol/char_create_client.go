package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type CharCreateRequest struct {
	Name        string
	Race        uint8
	Class       uint8
	Gender      uint8
	Skin        uint8
	Face        uint8
	HairStyle   uint8
	HairColor   uint8
	FacialStyle uint8
}

type charCreateWire struct {
	Race        uint8
	Class       uint8
	Gender      uint8
	Skin        uint8
	Face        uint8
	HairStyle   uint8
	HairColor   uint8
	FacialStyle uint8
}

func (c CharCreateRequest) isClientMessage() {}

func (c CharCreateRequest) Opcode() ClientOpcode {
	return ClientOpcodeCharCreate
}

func DecodeCharCreate(size uint16, reader io.Reader) (CharCreateRequest, error) {
	packetBody, err := readCharCreatePacket(size, reader)
	if err != nil {
		return CharCreateRequest{}, fmt.Errorf("error reading packet body: %w", err)
	}
	request, err := decodeCharCreateRequest(packetBody)
	if err != nil {
		return CharCreateRequest{}, fmt.Errorf("error decoding packet into request: %w", err)
	}
	return request, nil
}

func decodeCharCreateRequest(packetBody []byte) (CharCreateRequest, error) {
	reader := bytes.NewReader(packetBody)
	nameBuffer := make([]byte, len(packetBody)-9)
	_, err := io.ReadFull(reader, nameBuffer)
	if err != nil {
		return CharCreateRequest{}, fmt.Errorf("error reading character name from packet: %w", err)
	}
	var wire charCreateWire
	err = binary.Read(reader, binary.LittleEndian, &wire)
	if err != nil {
		return CharCreateRequest{}, fmt.Errorf("error reading bytes for charCreateWire: %w", err)
	}
	return CharCreateRequest{
		Name:        string(nameBuffer[:len(nameBuffer)-1]),
		Race:        wire.Race,
		Class:       wire.Class,
		Gender:      wire.Gender,
		Skin:        wire.Skin,
		Face:        wire.Face,
		HairStyle:   wire.HairStyle,
		HairColor:   wire.HairColor,
		FacialStyle: wire.FacialStyle,
	}, nil
}

func readCharCreatePacket(size uint16, reader io.Reader) ([]byte, error) {
	bodySize := size - ClientOpcodeLength
	body := make([]byte, bodySize)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
