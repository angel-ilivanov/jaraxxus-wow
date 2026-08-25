package protocol

import (
	"bytes"
	"encoding/binary"
)

var (
	size = []byte{0x00, 0x06} // Big endian, 6
)

type AuthChallengeServerMessage struct {
	ServerSeed []byte
}

func (a AuthChallengeServerMessage) Encode() []byte {
	var buf bytes.Buffer
	buf.Write(size)

	opcode := make([]byte, 2)
	binary.LittleEndian.AppendUint16(opcode, uint16(OpcodeAuthChallenge))
	buf.Write(opcode)
	buf.Write(a.ServerSeed)
	return buf.Bytes()
}
