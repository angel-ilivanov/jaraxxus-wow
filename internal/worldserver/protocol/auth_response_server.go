package protocol

import (
	"bytes"
	"encoding/binary"
)

type AuthResponse struct {
	ResultCode AccountResultValue
}

var billingPadding = make([]byte, 19)

const expansionCode uint8 = 2

func (a AuthResponse) Opcode() ServerOpcode {
	return OpcodeAuthResponse
}

func (a AuthResponse) EncodeBody() []byte {
	var packet bytes.Buffer

	resultBytes := make([]byte, 0, 4)
	resultBytes = binary.LittleEndian.AppendUint32(resultBytes, uint32(a.ResultCode))
	packet.Write(resultBytes)
	packet.Write(billingPadding)
	packet.WriteByte(expansionCode)

	return packet.Bytes()
}
