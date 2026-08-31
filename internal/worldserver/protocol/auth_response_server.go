package protocol

import (
	"bytes"
)

type AuthResponse struct {
	ResultCode AccountResultValue
}

var billingPadding = make([]byte, 9)

const expansionCode uint8 = 2

func (a AuthResponse) Opcode() ServerOpcode {
	return ServerOpcodeAuthResponse
}

func (a AuthResponse) EncodeBody() []byte {
	var packet bytes.Buffer

	packet.WriteByte(byte(a.ResultCode))
	if a.ResultCode != ResultAuthOk {
		return packet.Bytes()
	}

	packet.Write(billingPadding)
	packet.WriteByte(expansionCode)

	return packet.Bytes()
}
