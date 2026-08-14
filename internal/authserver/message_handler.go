package protocol

import (
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver"
)

func HandleMessage(session authserver.AuthSession, message ClientMessage) error {
	switch message.Opcode() {
	case 0x00:
	default:
		return fmt.Errorf("unknown opcode")
	}
	return nil
}
