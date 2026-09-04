package protocol

import "bytes"

type LoginVerifyWorldResponse struct {
	MapID     uint32
	PositionX float32
	PositionY float32
	PositionZ float32
}

func (l LoginVerifyWorldResponse) EncodeBody() []byte {
	var buf bytes.Buffer
	writeUint32LE(&buf, l.MapID)
	writeFloat32LE(&buf, l.PositionX)
	writeFloat32LE(&buf, l.PositionY)
	writeFloat32LE(&buf, l.PositionZ)
	return buf.Bytes()
}

func (l LoginVerifyWorldResponse) Opcode() ServerOpcode {
	return ServerOpcodeLoginVerifyWorld
}
