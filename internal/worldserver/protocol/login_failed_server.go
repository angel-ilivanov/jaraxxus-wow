package protocol

type LoginFailedResponse struct {
	Result AccountResultValue
}

func (l LoginFailedResponse) EncodeBody() []byte {
	return []byte{byte(l.Result)}
}

func (l LoginFailedResponse) Opcode() ServerOpcode {
	return ServerOpcodeCharLoginFailed
}
