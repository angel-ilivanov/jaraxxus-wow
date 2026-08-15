package authserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

type MessageHandler struct {
	store *accountstore.Store
}

// HandleMessage dispatches to the appropriate handler
func (handler *MessageHandler) HandleMessage(ctx context.Context, session *AuthSession, message protocol.ClientMessage) error {
	switch message.Opcode() {
	case 0x00:
		handler.handleLogonChallengeMessage(ctx, session, message.(protocol.LogonChallengeClientMessage))
	default:
		return fmt.Errorf("unknown opcode")
	}
	return nil
}

func NewHandler(store *accountstore.Store) *MessageHandler {
	return &MessageHandler{store: store}
}
