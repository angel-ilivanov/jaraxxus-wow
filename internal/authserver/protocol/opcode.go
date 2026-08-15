package protocol

type Opcode uint8

const (
	CMD_AUTH_LOGON_CHALLENGE Opcode = 0x00
	CMD_AUTH_LOGON_PROOF     Opcode = 0x01
)
