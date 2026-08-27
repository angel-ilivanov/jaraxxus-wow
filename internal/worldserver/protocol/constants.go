package protocol

type ServerOpcode uint16

const (
	OpcodeAuthChallenge ServerOpcode = 0x1EC
	OpcodeAuthResponse  ServerOpcode = 0x1EE
)

type ClientOpcode uint32

const (
	OpcodeAuthSession ClientOpcode = 0x1ED
	OpcodeCharEnum    ClientOpcode = 0x037
)

const (
	MaxUsernameLength  = 16
	ClientOpcodeLength = 4
	ServerOpcodeLength = 2
)

type AccountResultValue uint32

const (
	ResultAuthOk     AccountResultValue = 0x0C
	ResultAuthReject AccountResultValue = 0x0E
)
