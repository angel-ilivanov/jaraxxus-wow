package worldserver

import (
	"context"
	"errors"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

var ErrUnknownRequestType = errors.New("unknown request type")

func HandleRequest(ctx context.Context, session *session, request protocol.ClientMessage) (protocol.ServerMessage, error) {
	switch request := request.(type) {
	case protocol.AuthSessionRequest:
		return handleAuthSession(ctx, session, request)
	default:
		return nil, ErrUnknownRequestType
	}
}
