package auth

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type wireCmdAuthLogonChallengeClientHeader struct {
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
type CmdAuthLogonChallengeClient struct {
	opcode      uint8
	ip          uint32
	accountName string
}

// 00082100576f57000303053430363878006e69570053556e6520feffffc0a8728203414755s
func ParsePacket(data []byte) (wireCmdAuthLogonChallengeClientHeader, []byte) {
	var h wireCmdAuthLogonChallengeClientHeader
	reader := bytes.NewReader(data)

	err := binary.Read(reader, binary.LittleEndian, &h)
	if err != nil {
		fmt.Printf("Error parsing binary: %s", err)
	}

	nameBytes := make([]byte, h.AccountNameLength)
	_, err = io.ReadFull(reader, nameBytes)
	if err != nil {
		return h, nameBytes
	}
	parseInfo(h, nameBytes)
	return h, nameBytes
}

func parseInfo(header wireCmdAuthLogonChallengeClientHeader, nameBytes []byte) CmdAuthLogonChallengeClient {
	ip := binary.BigEndian.Uint32(header.Ip[:])
	parsed := CmdAuthLogonChallengeClient{
		opcode:      header.Opcode,
		ip:          ip,
		accountName: string(nameBytes),
	}

	fmt.Printf("Opcode: %d\n", parsed.opcode)
	fmt.Printf("IP: %d\n", parsed.ip)
	fmt.Printf("Account Name: %s\n", parsed.accountName)

	return parsed

}
