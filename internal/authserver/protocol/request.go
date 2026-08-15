package protocol

type Request interface {
	Opcode() Opcode
}
