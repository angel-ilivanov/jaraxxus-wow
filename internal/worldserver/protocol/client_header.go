package protocol

type ClientHeader struct {
	PacketSize uint16 // Size of the remaining packet (opcode + body)
	Opcode     ClientOpcode
}
