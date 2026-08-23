package protocol

type Request interface {
	isRequest()
	Opcode() Opcode
}

type Response interface {
	isResponse()
	Opcode() Opcode
}
