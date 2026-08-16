package protocol

import (
	"bytes"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

type LogonChallengeResponse struct {
	Result          LoginResult
	ServerPublicKey []byte
	Salt            []byte
}

var crcSalt = make([]byte, 16)

func EncodeLogonChallengeResponse(response LogonChallengeResponse) []byte {
	var buf bytes.Buffer
	// Opcode
	buf.WriteByte(byte(CmdAuthLogonChallenge))
	// Protocol Version: 0
	buf.WriteByte(0x00)
	buf.WriteByte(byte(response.Result))

	if response.Result != ResultSuccess {
		fmt.Println("Server Packet:", buf.Bytes())
		return buf.Bytes()
	}

	buf.Write(response.ServerPublicKey)
	buf.WriteByte(srp6.GeneratorLength)
	buf.WriteByte(srp6.Generator)
	buf.WriteByte(srp6.LargeSafePrimeLength)
	buf.Write(srp6.LargeSafePrimeLittleEndian)
	buf.Write(response.Salt)
	// CRC Salt
	buf.Write(crcSalt)
	// Additional Verification required (PIN, 2FA)
	buf.WriteByte(0x00)
	fmt.Println("Server Packet:", buf.Bytes())
	return buf.Bytes()
}
