package worldserver

import (
	"context"
	"errors"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handlePlayerLogin(ctx context.Context, session *session, request protocol.PlayerLoginRequest) (protocol.ServerMessage, error) {
	valid, err := handler.characterBelongsToAccount(ctx, session, request)

	if err != nil || !valid {
		return protocol.LoginFailedResponse{Result: protocol.ResultCharLoginFailed}, err
	}

	spawnPoint, err := handler.characterStore.FindCharacterSpawnPoint(ctx, request.CharGUID)
	if err != nil {
		return protocol.LoginFailedResponse{Result: protocol.ResultCharLoginFailed}, err
	}
	session.ActiveCharacterGuid = request.CharGUID
	return protocol.LoginVerifyWorldResponse{
		MapID:       spawnPoint.MapId,
		PositionX:   spawnPoint.PositionX,
		PositionY:   spawnPoint.PositionY,
		PositionZ:   spawnPoint.PositionZ,
		Orientation: spawnPoint.Orientation,
	}, nil
}

func (handler RequestHandler) characterBelongsToAccount(ctx context.Context, session *session, request protocol.PlayerLoginRequest) (bool, error) {
	_, err := handler.characterStore.GetByAccountIdAndGUID(ctx, session.AccountId, request.CharGUID)
	if err != nil {
		if errors.Is(err, characterstore.ErrCharacterNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
