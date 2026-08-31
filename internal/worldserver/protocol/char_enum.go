package protocol

import "bytes"

type CharEnumRequest struct{}

func (c CharEnumRequest) isClientMessage() {}

func (c CharEnumRequest) Opcode() ClientOpcode {
	return ClientOpcodeCharEnum
}

type CharEnumResponse struct {
	AmountOfCharacters uint8
}

func (c CharEnumResponse) EncodeBody() []byte {
	var buf bytes.Buffer
	buf.WriteByte(c.AmountOfCharacters)
	if c.AmountOfCharacters == 0 {
		return buf.Bytes()
	}
	// characters...
	return buf.Bytes()
}

func (c CharEnumResponse) Opcode() ServerOpcode {
	return ServerOpcodeCharEnum
}
