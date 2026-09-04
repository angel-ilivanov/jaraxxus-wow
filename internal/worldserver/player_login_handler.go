package worldserver

import (
	"context"
	"errors"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handlePlayerLogin(ctx context.Context, session *session, request protocol.PlayerLoginRequest) ([]protocol.ServerMessage, error) {
	valid, err := handler.characterBelongsToAccount(ctx, session, request)
	failResponse := protocol.LoginFailedResponse{Result: protocol.ResultCharLoginFailed}

	if err != nil || !valid {
		return []protocol.ServerMessage{failResponse}, err
	}

	spawnPoint, err := handler.characterStore.FindCharacterSpawnPoint(ctx, request.CharGUID)
	if err != nil {
		return []protocol.ServerMessage{failResponse}, err
	}
	session.ActiveCharacterGuid = request.CharGUID

	verifyWorldResponse := protocol.LoginVerifyWorldResponse{
		MapID:       spawnPoint.MapId,
		PositionX:   spawnPoint.PositionX,
		PositionY:   spawnPoint.PositionY,
		PositionZ:   spawnPoint.PositionZ,
		Orientation: spawnPoint.Orientation,
	}
	return []protocol.ServerMessage{verifyWorldResponse, protocol.TutorialFlagsResponse{}}, nil
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
