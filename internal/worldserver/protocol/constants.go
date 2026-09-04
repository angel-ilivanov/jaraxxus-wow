package protocol

type ServerOpcode uint16

const (
	ServerOpcodeAuthChallenge ServerOpcode = 0x1EC
	ServerOpcodeAuthResponse  ServerOpcode = 0x1EE
	ServerOpcodeAccountTimes  ServerOpcode = 0x209
	ServerOpcodeCharEnum      ServerOpcode = 0x03B
	ServerOpcodeCharCreate    ServerOpcode = 0x03A
)

type ClientOpcode uint32

const (
	ClientOpcodeAuthSession       ClientOpcode = 0x1ED
	ClientOpcodeCharEnum          ClientOpcode = 0x037
	ClientOpcodeAccountTimesReady ClientOpcode = 0x4FF
	ClientOpcodeCharCreate        ClientOpcode = 0x0036
	ClientOpcodePlayerLogin       ClientOpcode = 0x003D
)

const (
	MaxUsernameLength  = 16
	ClientOpcodeLength = 4
	ServerOpcodeLength = 2
)

type AccountResultValue uint8

const (
	ResultAuthOk              AccountResultValue = 0x0C
	ResultAuthReject          AccountResultValue = 0x0E
	ResultCharCreateSuccess   AccountResultValue = 0x2F
	ResultCharCreateNameInUse AccountResultValue = 0x32
)
