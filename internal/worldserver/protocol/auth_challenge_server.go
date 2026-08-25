package protocol

import (
	"bytes"
	"encoding/binary"
)

const bodySize uint16 = 42

var dosChallenge = make([]byte, 32)

type AuthChallengeServerMessage struct {
	ServerSeed []byte
}

func (a AuthChallengeServerMessage) Encode() []byte {
	var buf bytes.Buffer

	// Header:
	var size []byte
	size = binary.BigEndian.AppendUint16(size, bodySize)
	buf.Write(size)

	var opcode []byte
	opcode = binary.LittleEndian.AppendUint16(opcode, uint16(OpcodeAuthChallenge))
	buf.Write(opcode)

	// Body:
	var dosDifficulty []byte
	dosDifficulty = binary.LittleEndian.AppendUint32(dosDifficulty, 1)
	buf.Write(dosDifficulty)

	buf.Write(a.ServerSeed)
	buf.Write(dosChallenge)
	return buf.Bytes()
}
