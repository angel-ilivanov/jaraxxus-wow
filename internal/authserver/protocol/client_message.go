package protocol

type ClientMessage interface {
	Opcode() Opcode
}
