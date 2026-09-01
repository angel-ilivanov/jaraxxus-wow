package protocol

type ServerOpcode uint16

const (
	OpcodeAuthChallenge ServerOpcode = 0x1EC
)

type ClientOpcode uint32

const (
	OpcodeAuthSession ClientOpcode = 0x1ED
)
