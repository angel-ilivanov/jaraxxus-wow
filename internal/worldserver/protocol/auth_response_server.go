package protocol

import (
	"bytes"
	"encoding/binary"
)

type AuthResponse struct {
	ResultCode AccountResultValue
}

func (a AuthResponse) Opcode() ServerOpcode {
	return OpcodeAuthResponse
}

func (a AuthResponse) EncodeBody() []byte {
	var packet bytes.Buffer

	resultBytes := make([]byte, 0, 4)
	resultBytes = binary.LittleEndian.AppendUint32(resultBytes, uint32(a.ResultCode))
	packet.Write(resultBytes)

	return packet.Bytes()
}
