package worldserver

import (
	"context"
	"errors"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

var ErrUnknownRequestType = errors.New("unknown request type")

type RequestHandler struct {
	store *accountstore.Store
}

func (handler RequestHandler) HandleRequest(ctx context.Context, session *session, request protocol.ClientMessage) (protocol.ServerMessage, error) {
	switch request := request.(type) {
	case protocol.AuthSessionRequest:
		return handler.handleAuthSession(ctx, session, request)
	case protocol.AccountDataTimesRequest:
		return protocol.AccountDataTimesResponse{}, nil
	case protocol.CharEnumRequest:
		return handler.handleCharEnum(ctx, session, request)
	case protocol.CharCreateRequest:
		return handler.handleCharCreate(ctx, session, request)
		//return nil, fmt.Errorf("handling char creation not yet implemented")
	default:
		return nil, ErrUnknownRequestType
	}
}

func NewRequestHandler(store *accountstore.Store) *RequestHandler {
	return &RequestHandler{store: store}
}
