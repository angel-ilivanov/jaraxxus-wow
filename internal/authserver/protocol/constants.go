package protocol

type LoginResult uint8

const (
	ResultSuccess        LoginResult = 0x00
	ResultUnknownAccount             = 0x04
)

type Opcode uint8

const (
	CmdAuthLogonChallenge Opcode = 0x00
	CmdAuthLogonProof            = 0x01
	CmdRealmList                 = 0x10
)
