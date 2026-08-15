package protocol

type LoginPacketResult uint8

const (
	SUCCESS              LoginPacketResult = 0x00
	FAIL_UNKNOWN_ACCOUNT                   = 0x01
)
