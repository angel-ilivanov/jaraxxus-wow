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
func (handler *MessageHandler) HandleMessage(ctx context.Context, session *AuthSession, message protocol.ClientMessage) ([]byte, error) {
	switch message := message.(type) {
	case protocol.LogonChallengeClientMessage:
		return handler.handleLogonChallengeMessage(ctx, session, message)
	default:
		return nil, fmt.Errorf("unknown opcode")
	}
}

func NewHandler(store *accountstore.Store) *MessageHandler {
	return &MessageHandler{store: store}
}
