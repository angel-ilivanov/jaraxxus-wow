package protocol

type ServerOpcode uint16

const (
	OpcodeAuthChallenge ServerOpcode = 0x1EC
)

type ClientOpcode uint32

const (
	OpcodeAuthSession ClientOpcode = 0x1ED
)

const (
	MaxUsernameLength  = 16
	ClientOpcodeLength = 4
	ServerOpcodeLength = 2
)

type AccountResultValue uint32

const (
	ResultSuccess AccountResultValue = 0x00
	ResultFailure AccountResultValue = 0x01
)
