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
	var header bytes.Buffer
	sizeBytes := make([]byte, 0, 2)
	sizeBytes = binary.BigEndian.AppendUint16(sizeBytes, size)
	header.Write(sizeBytes)

	opcodeBytes := make([]byte, 0, 2)
	opcodeBytes = binary.LittleEndian.AppendUint16(opcodeBytes, uint16(OpcodeAuthResponse))
	header.Write(opcodeBytes)

	var body bytes.Buffer
	resultBytes := make([]byte, 0, 4)
	resultBytes = binary.LittleEndian.AppendUint32(resultBytes, uint32(a.ResultCode))
	body.Write(resultBytes)

	//encrypt header, return encrypted header | body
	return nil
}
