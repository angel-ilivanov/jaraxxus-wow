package protocol

import (
	"bytes"
	"encoding/binary"
)

var size = uint16(23 + len(realmName) + len(addressPort))
var headerPadding = make([]byte, 4)
var realmName = []byte("Jaraxxus\x00")
var addressPort = make([]byte, 2)
var population = []byte{0x00, 0x00, 0x00, 0x00}               //MEDIUM
var buildSpecification = []byte{0x03, 0x03, 0x05, 0x34, 0x30} //3.3.5, build 12340
var footerPadding = make([]byte, 2)

type RealmListResponse struct {
	NumChars uint8
}

func (r RealmListResponse) isResponse() {}

func EncodeRealmListResponse(response RealmListResponse) []byte {
	sizeLittleEndian := make([]byte, 2)
	binary.LittleEndian.PutUint16(sizeLittleEndian, size)

	var buf bytes.Buffer
	buf.WriteByte(CmdRealmList)
	//Size
	buf.Write(sizeLittleEndian)
	buf.Write(headerPadding)
	buf.Write([]byte{0x01, 0x00}) //NUMBER OF REALMS, TODO: SHOW NICER REPRESENTATION
	//realm type:
	buf.WriteByte(0x00)
	//locked:
	buf.WriteByte(0x00)
	//realm flags:
	buf.WriteByte(0x04)
	buf.Write(realmName)
	buf.Write(addressPort) //2 bytes currently
	buf.Write(population)
	buf.WriteByte(response.NumChars)
	//realm category:
	buf.WriteByte(0x00) //DEFAULT
	//realm id:
	buf.WriteByte(0x00)
	buf.Write(buildSpecification)
	buf.Write(footerPadding)
	return buf.Bytes()
}
