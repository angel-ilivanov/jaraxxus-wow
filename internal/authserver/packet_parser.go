package authserver

import (
	"log"

	"github.com/angel-ilivanov/wow-server/internal/authserver/protocol"
)

func Parse(packet []byte) {
	if len(packet) == 0 {
		log.Fatal("Insufficient packet size")
	}
	opcode := packet[0]

	switch opcode {
	case 0:
		protocol.ParsePacket(packet)
	}
}
