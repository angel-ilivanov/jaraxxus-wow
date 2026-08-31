package worldserver

import (
	"context"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handleCharEnum(ctx context.Context, session *session, request protocol.CharEnumRequest) (protocol.CharEnumResponse, error) {
	// TODO: Fetch account characters
	return protocol.CharEnumResponse{AmountOfCharacters: 0}, nil
}
