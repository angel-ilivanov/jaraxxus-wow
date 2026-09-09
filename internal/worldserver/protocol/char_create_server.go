package protocol

type CharCreateResponse struct {
	Result AccountResultValue
}

func (c CharCreateResponse) EncodeBody() []byte {
	return []byte{byte(c.Result)}
}

func (c CharCreateResponse) Opcode() ServerOpcode {
	return ServerOpcodeCharCreate
}
