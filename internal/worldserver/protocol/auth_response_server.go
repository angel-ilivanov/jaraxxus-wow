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
	return OpcodeAuthResponse
}

func (a AuthResponse) EncodeBody() []byte {
	var packet bytes.Buffer

	packet.WriteByte(byte(a.ResultCode))
	packet.Write(billingPadding)
	packet.WriteByte(expansionCode)

	return packet.Bytes()
}
