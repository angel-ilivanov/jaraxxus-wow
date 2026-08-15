package protocol

import (
	"bytes"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

type LogonChallengeServerMessage struct {
	Success         bool
	ServerPublicKey []byte
	Salt            []byte
}

var crcSalt = make([]byte, 16)

func AssembleServerChallengePacket(message LogonChallengeServerMessage) []byte {
	var buf bytes.Buffer
	// Opcode: CMD_AUTH_LOGON_CHALLENGE
	buf.WriteByte(0x00)
	// Protocol Version: 0
	buf.WriteByte(0x00)
	// Result: SUCCESS (0)
	buf.WriteByte(0x00)
	buf.Write(message.ServerPublicKey)
	buf.WriteByte(srp6.GeneratorLength)
	buf.WriteByte(srp6.Generator)
	buf.WriteByte(srp6.LargeSafePrimeLength)
	buf.Write(srp6.LargeSafePrimeLittleEndian)
	buf.Write(message.Salt)
	// CRC Salt
	buf.Write(crcSalt)
	// Additional Verification required (PIN, 2FA)
	buf.WriteByte(0x00)
	fmt.Println("Server Packet:", buf.Bytes())
	return buf.Bytes()
}
