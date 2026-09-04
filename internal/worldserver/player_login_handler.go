package worldserver

import (
	"context"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handlePlayerLogin(ctx context.Context, session *session, request protocol.PlayerLoginRequest) (protocol.ServerMessage, error) {
	valid, err := handler.characterBelongsToAccount(ctx, session, request)

	// TODO: CHECK VALIDITY
	if err != nil || !valid {
		return protocol.LoginFailedResponse{Result: protocol.ResultCharLoginFailed}, err
	}

	// TODO FETCH FROM DB:
	return protocol.LoginVerifyWorldResponse{
		MapID:     0,
		PositionX: 0,
		PositionY: 0,
		PositionZ: 0,
	}, nil
}

func (handler RequestHandler) characterBelongsToAccount(ctx context.Context, session *session, request protocol.PlayerLoginRequest) (bool, error) {
	return true, nil
}
