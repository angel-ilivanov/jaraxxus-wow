package protocol

type ServerMessage interface {
	EncodeBody() []byte
	Opcode() ServerOpcode
}

type ClientMessage interface {
	isClientMessage()
	Opcode() ClientOpcode
}
