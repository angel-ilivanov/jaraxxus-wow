package worldserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

var ErrUnknownRequestType = errors.New("unknown request type")

type RequestHandler struct {
	accountStore   *accountstore.Store
	characterStore *characterstore.Store
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
	case protocol.PlayerLoginRequest:
		return nil, fmt.Errorf("handling player login request not implemented yet")
	default:
		return nil, ErrUnknownRequestType
	}
}

func NewRequestHandler(accountStore *accountstore.Store, characterStore *characterstore.Store) *RequestHandler {
	return &RequestHandler{accountStore: accountStore, characterStore: characterStore}
}
