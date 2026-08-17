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

// HandleRequest dispatches to the appropriate handler
func (handler *RequestHandler) HandleRequest(ctx context.Context, session *AuthSession, request protocol.Request) (protocol.Response, error) {
	switch request := request.(type) {
	case protocol.LogonChallengeRequest:
		return handler.handleLogonChallenge(ctx, session, request)
	case protocol.LogonProofRequest:
		return handler.handleLogonProof(ctx, session, request)
	case protocol.RealmListRequest:
		return nil, fmt.Errorf("responding to realm list request not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported request type %T", request)
	}
}

func NewRequestHandler(store *accountstore.Store) *RequestHandler {
	return &RequestHandler{store: store}
}
