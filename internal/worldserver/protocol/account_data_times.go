package protocol

import (
	"bytes"
	"encoding/binary"
	"time"
)

type AccountDataTimesRequest struct {
}

func (a AccountDataTimesRequest) isClientMessage() {
}

func (a AccountDataTimesRequest) Opcode() ClientOpcode {
	return ClientOpcodeAccountTimesReady
}

type AccountDataTimesResponse struct{}

func (a AccountDataTimesResponse) EncodeBody() []byte {
	var buf bytes.Buffer

	timeBuffer := make([]byte, 0, 4)
	currentUnixTime := uint32(time.Now().Unix())
	binary.LittleEndian.AppendUint32(timeBuffer, currentUnixTime)
	buf.Write(timeBuffer)

	buf.WriteByte(1) // Unknown value hardcoded to 1

	settingsMaskBuffer := make([]byte, 0, 4)
	settingsMask := uint32(0x15)
	binary.LittleEndian.AppendUint32(settingsMaskBuffer, settingsMask)
	buf.Write(timeBuffer)

	timestampPadding := make([]byte, 12)
	buf.Write(timestampPadding)
	return buf.Bytes()
}

func (a AccountDataTimesResponse) Opcode() ServerOpcode {
	return ServerOpcodeAccountTimes
}
