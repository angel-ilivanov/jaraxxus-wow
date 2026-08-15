package protocol

type LoginResult uint8

const (
	Success            LoginResult = 0x00
	FailUnknownAccount             = 0x01
)

type Opcode uint8

const (
	CmdAuthLogonChallenge Opcode = 0x00
	CmdAuthLogonProof     Opcode = 0x01
)
