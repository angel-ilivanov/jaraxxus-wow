package worldserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handleCharEnum(ctx context.Context, session *session, request protocol.CharEnumRequest) (protocol.CharEnumResponse, error) {
	// TODO: Fetch account characters
	characters, err := handler.characterStore.ListByAccountId(ctx, session.AccountId)
	if err != nil {
		return protocol.CharEnumResponse{}, fmt.Errorf("error retrieving characters: %w", err)
	}
	return protocol.CharEnumResponse{Characters: characters}, nil
}
