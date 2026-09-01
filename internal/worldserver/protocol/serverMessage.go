package protocol

type serverMessage interface {
	Encode() []byte
}

type ClientMessage interface {
	isClientMessage()
	Opcode() ClientOpcode
}
