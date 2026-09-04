package worldserver

import (
	"context"
	"errors"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

var ErrUnknownRequestType = errors.New("unknown request type")

type RequestHandler struct {
	accountStore   *accountstore.Store
	characterStore *characterstore.Store
}

func (handler RequestHandler) HandleRequest(ctx context.Context, session *session, request protocol.ClientMessage) ([]protocol.ServerMessage, error) {
	switch request := request.(type) {
	case protocol.AuthSessionRequest:
		response, err := handler.handleAuthSession(ctx, session, request)
		return []protocol.ServerMessage{response}, err
	case protocol.AccountDataTimesRequest:
		return []protocol.ServerMessage{protocol.AccountDataTimesResponse{}}, nil
	case protocol.CharEnumRequest:
		response, err := handler.handleCharEnum(ctx, session, request)
		return []protocol.ServerMessage{response}, err
	case protocol.CharCreateRequest:
		response, err := handler.handleCharCreate(ctx, session, request)
		return []protocol.ServerMessage{response}, err
	case protocol.PlayerLoginRequest:
		return handler.handlePlayerLogin(ctx, session, request)
	default:
		return nil, ErrUnknownRequestType
	}
}

func NewRequestHandler(accountStore *accountstore.Store, characterStore *characterstore.Store) *RequestHandler {
	return &RequestHandler{accountStore: accountStore, characterStore: characterStore}
}
