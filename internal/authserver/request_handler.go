package authserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

type RequestHandler struct {
	store *accountstore.Store
}

// HandleMessage dispatches to the appropriate handler
func (handler *RequestHandler) HandleMessage(ctx context.Context, session *AuthSession, request protocol.Request) ([]byte, error) {
	switch message := request.(type) {
	case protocol.LogonChallengeRequest:
		return handler.handleLogonChallenge(ctx, session, message)
	default:
		return nil, fmt.Errorf("unknown opcode")
	}
}

func NewRequestHandler(store *accountstore.Store) *RequestHandler {
	return &RequestHandler{store: store}
}
