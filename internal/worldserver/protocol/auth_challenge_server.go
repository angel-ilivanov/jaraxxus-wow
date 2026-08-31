package protocol

import (
	"bytes"
	"encoding/binary"
)

var dosChallenge = make([]byte, 32) // unused DoS challenge information, treat as padding

type AuthChallengeServerMessage struct {
	ServerSeed []byte
}

func (a AuthChallengeServerMessage) Opcode() ServerOpcode {
	return ServerOpcodeAuthChallenge
}

func (a AuthChallengeServerMessage) EncodeBody() []byte {
	var buf bytes.Buffer

	var dosDifficulty []byte // unused, typically hardcoded to 1
	dosDifficulty = binary.LittleEndian.AppendUint32(dosDifficulty, 1)
	buf.Write(dosDifficulty)

	buf.Write(a.ServerSeed)
	buf.Write(dosChallenge)
	return buf.Bytes()
}
