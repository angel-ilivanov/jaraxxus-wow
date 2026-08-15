package authserver

import (
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

type MessageHandler struct {
	store *accountstore.Store
}

func (handler *MessageHandler) HandleMessage(session *AuthSession, message protocol.ClientMessage) error {
	switch message.Opcode() {
	case 0x00:

	default:
		return fmt.Errorf("unknown opcode")
	}
	return nil
}

func NewHandler(store *accountstore.Store) *MessageHandler {
	return &MessageHandler{store: store}
}
