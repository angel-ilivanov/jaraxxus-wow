package worldserver

import (
	"context"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handleCharCreate(ctx context.Context, session *session, request protocol.CharCreateRequest) (protocol.CharCreateResponse, error) {
	// TODO: create character in characterStore
	return protocol.CharCreateResponse{Result: protocol.ResultCharCreateSuccess}, nil
}
