package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Read and decode the challenge-specific layout

type logonChallengeWire struct {
	Opcode            uint8
	ProtocolVersion   uint8
	Size              uint16
	GameName          [4]uint8
	Version           [3]uint8
	Build             uint16
	Platform          [4]uint8
	Os                [4]uint8
	Locale            [4]uint8 //reversed but irrelevant
	WorldRegionBias   uint32
	Ip                [4]uint8
	AccountNameLength uint8
}

type LogonChallengeRequest struct {
	Ip          uint32
	AccountName string
}

func (l LogonChallengeRequest) Opcode() Opcode {
	return CmdAuthLogonChallenge
}

func DecodeLogonChallengeRequest(conn io.Reader) (Request, error) {
	packetBytes, err := assemblePacket(conn)
	if err != nil {
		return LogonChallengeRequest{}, fmt.Errorf("error assembling logon challenge Packet %v", err)
	}
	parsedPacket, err := mapToStruct(packetBytes)

	if err != nil {
		return LogonChallengeRequest{}, fmt.Errorf("error mapping logon challenge Packet to struct %v", err)
	}
	return constructRequest(parsedPacket, packetBytes), nil
}
func assemblePacket(conn io.Reader) ([]byte, error) {
	header := make([]byte, 3)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}
	header = append([]byte{0x00}, header...) // prepend opcode
	fmt.Println("header:", header)
	size := binary.LittleEndian.Uint16(header[2:4])
	body := make([]byte, size)
	_, err = io.ReadFull(conn, body)
	if err != nil {
		return nil, err
	}
	fmt.Println("body:", body)
	return append(header, body...), nil
}

func mapToStruct(packetBytes []byte) (logonChallengeWire, error) {
	var parsed logonChallengeWire
	reader := bytes.NewReader(packetBytes)
	err := binary.Read(reader, binary.LittleEndian, &parsed)
	if err != nil {
		return logonChallengeWire{}, err
	}
	fmt.Println("Account name length:", parsed.AccountNameLength)
	return parsed, nil
}

func readName(parsedPacket logonChallengeWire, packetBytes []byte) string {
	nameLength := int(parsedPacket.AccountNameLength)
	nameBytes := packetBytes[len(packetBytes)-nameLength:]
	fmt.Println("Account name:", string(nameBytes))
	return string(nameBytes)
}

func constructRequest(parsedPacket logonChallengeWire, packetBytes []byte) LogonChallengeRequest {
	return LogonChallengeRequest{
		Ip:          binary.BigEndian.Uint32(parsedPacket.Ip[:]),
		AccountName: readName(parsedPacket, packetBytes),
	}
}
