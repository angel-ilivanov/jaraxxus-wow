package protocol

type ServerMessage interface {
	Encode() []byte
}

type ClientMessage interface {
	isClientMessage()
	Opcode() ClientOpcode
}
