package protocol

import (
	"bytes"
)

type LogonProofResponse struct {
	Result      LoginResult
	ServerProof []byte
}

func (l LogonProofResponse) Opcode() Opcode {
	return CmdAuthLogonProof
}

func (l LogonProofResponse) isResponse() {}

// 4 bytes if unsuccessful
// 32 bytes if successful
func EncodeLogonProofResponse(response LogonProofResponse) []byte {
	var buf bytes.Buffer
	// Opcode
	buf.WriteByte(byte(CmdAuthLogonProof))
	buf.WriteByte(byte(response.Result))
	if response.Result != ResultSuccess {
		buf.Write([]byte{0, 0}) //padding
		return buf.Bytes()
	}
	buf.Write(response.ServerProof)
	buf.Write([]byte{00, 00, 80, 00}) //ACCOUNT_FLAG_PROPASS = 0x00800000, mangos hardcodes this
	buf.Write([]byte{00, 00, 00, 00}) // hardware survey, set to 0
	buf.Write([]byte{00, 00})         // unknown flags, set to 0
	return buf.Bytes()
}
