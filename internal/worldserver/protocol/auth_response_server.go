package protocol

import (
	"bytes"
	"encoding/binary"
)

const size = 6

type AuthResponse struct {
	ResultCode AccountResultValue
}

func (a AuthResponse) Encode() []byte {
	var packet bytes.Buffer

	sizeBytes := make([]byte, 0, 2)
	sizeBytes = binary.BigEndian.AppendUint16(sizeBytes, size)
	packet.Write(sizeBytes)

	opcodeBytes := make([]byte, 0, ServerOpcodeLength)
	opcodeBytes = binary.LittleEndian.AppendUint16(opcodeBytes, uint16(OpcodeAuthResponse))
	packet.Write(opcodeBytes)

	resultBytes := make([]byte, 0, 4)
	resultBytes = binary.LittleEndian.AppendUint32(resultBytes, uint32(a.ResultCode))
	packet.Write(resultBytes)

	return packet.Bytes()
}
