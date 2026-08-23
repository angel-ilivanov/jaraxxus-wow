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
	OS                [4]uint8
	Locale            [4]uint8 //reversed but irrelevant
	WorldRegionBias   uint32
	IP                [4]uint8
	AccountNameLength uint8
}

type LogonChallengeRequest struct {
	IP          uint32
	AccountName string
}

func (l LogonChallengeRequest) isRequest() {}

func DecodeLogonChallengeRequest(reader io.Reader) (Request, error) {
	packetBytes, err := readLogonChallengePacket(reader)
	if err != nil {
		return LogonChallengeRequest{}, fmt.Errorf("error assembling logon challenge Packet %v", err)
	}
	parsedPacket, err := decodeLogonChallengeWire(packetBytes)

	if err != nil {
		return LogonChallengeRequest{}, fmt.Errorf("error mapping logon challenge Packet to struct %v", err)
	}
	return newLogonChallengeRequest(parsedPacket, packetBytes), nil
}
func readLogonChallengePacket(reader io.Reader) ([]byte, error) {
	header := make([]byte, 3)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		return nil, err
	}
	header = append([]byte{0x00}, header...) // prepend opcode
	size := binary.LittleEndian.Uint16(header[2:4])
	body := make([]byte, size)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func decodeLogonChallengeWire(packetBytes []byte) (logonChallengeWire, error) {
	var parsed logonChallengeWire
	reader := bytes.NewReader(packetBytes)
	err := binary.Read(reader, binary.LittleEndian, &parsed)
	if err != nil {
		return logonChallengeWire{}, err
	}
	return parsed, nil
}

func decodeAccountName(parsedPacket logonChallengeWire, packetBytes []byte) string {
	nameLength := int(parsedPacket.AccountNameLength)
	nameBytes := packetBytes[len(packetBytes)-nameLength:]
	return string(nameBytes)
}

func newLogonChallengeRequest(parsedPacket logonChallengeWire, packetBytes []byte) LogonChallengeRequest {
	return LogonChallengeRequest{
		IP:          binary.BigEndian.Uint32(parsedPacket.IP[:]),
		AccountName: decodeAccountName(parsedPacket, packetBytes),
	}
}
