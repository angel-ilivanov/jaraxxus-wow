package protocol

type ClientMessage interface {
	Opcode() uint8
}
