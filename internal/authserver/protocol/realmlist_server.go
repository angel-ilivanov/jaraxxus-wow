package protocol

import (
	"bytes"
	"encoding/binary"
)

var (
	headerPadding  = make([]byte, 4)
	population     = []byte{0x00, 0x00, 0x00, 0x00} //MEDIUM
	footerPadding  = make([]byte, 2)
	numberOfRealms = []byte{0x01, 0x00} // one (little endian)
)

const (
	realmName             = "Jaraxxus\x00"
	flagRecommendedRealm  = 0x20
	categoryEuropeEnglish = 0x01
	realmId               = 0x01
	realmTypePvp          = 0x01
	realmLockDisabled     = 0x00
)

type RealmListResponse struct {
	NumChars           uint8
	WorldServerAddress string
}

func (r RealmListResponse) isResponse() {}

func EncodeRealmListResponse(response RealmListResponse) []byte {
	addressPort := append([]byte(
		response.WorldServerAddress),
		0x00) // Append \0x00 byte
	var size = uint16(
		18 +
			len(realmName) +
			len(addressPort)) // packet size without opcode and size
	packetSizeLittleEndian := make([]byte, 2)
	binary.LittleEndian.PutUint16(packetSizeLittleEndian, size)

	var buf bytes.Buffer
	buf.WriteByte(CmdRealmList)
	buf.Write(packetSizeLittleEndian)
	buf.Write(headerPadding)
	buf.Write(numberOfRealms)
	buf.WriteByte(realmTypePvp)
	buf.WriteByte(realmLockDisabled)
	buf.WriteByte(flagRecommendedRealm)
	buf.WriteString(realmName)
	buf.Write(addressPort)
	buf.Write(population)
	buf.WriteByte(response.NumChars)
	buf.WriteByte(categoryEuropeEnglish)
	buf.WriteByte(realmId)
	buf.Write(footerPadding)
	return buf.Bytes()
}
