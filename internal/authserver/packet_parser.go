package authserver

import (
	"fmt"

	"github.com/angel-ilivanov/wow-server/internal/authserver/protocol"
)

func Parse(packet []byte) {
	if len(packet) == 0 {
		fmt.Errorf("Insufficient packet size")
	}
	opcode := packet[0]

	switch opcode {
	case 0:
		protocol.ParsePacket(packet)
	}
}
