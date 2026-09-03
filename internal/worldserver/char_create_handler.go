package worldserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handleCharCreate(ctx context.Context, session *session, request protocol.CharCreateRequest) (protocol.CharCreateResponse, error) {
	// TODO: create character in characterStore
	_, err := handler.characterStore.Insert(ctx, session.AccountId, generateCharacterFromRequest(request))
	if err != nil {
		fmt.Println(err)
		return protocol.CharCreateResponse{Result: protocol.ResultCharCreateNameInUse}, nil
	}
	return protocol.CharCreateResponse{Result: protocol.ResultCharCreateSuccess}, nil
}

func generateCharacterFromRequest(request protocol.CharCreateRequest) characterstore.NewCharacter {
	return characterstore.NewCharacter{
		Name:   request.Name,
		Race:   characterstore.Race(request.Race),
		Class:  characterstore.Class(request.Class),
		Gender: characterstore.Gender(request.Gender),
		Appearance: characterstore.Appearance{
			Skin:        request.Skin,
			Face:        request.Face,
			HairStyle:   request.HairStyle,
			HairColor:   request.HairColor,
			FacialStyle: request.FacialStyle,
		},
		//TODO FIX PLACEHOLDER VALUES:
		GeneratedState: characterstore.InitialState{
			Level:       1,
			MapID:       0,
			ZoneID:      0,
			PositionX:   0,
			PositionY:   0,
			PositionZ:   0,
			Orientation: 0,
		},
	}
}
