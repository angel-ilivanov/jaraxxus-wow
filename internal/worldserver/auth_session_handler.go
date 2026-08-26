package worldserver

import (
	"context"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func handleAuthSession(ctx context.Context, session *session, request protocol.AuthSessionRequest) (protocol.ServerMessage, error) {
	session.clientSeed = request.ClientSeed
	session.Username = request.Username
	return nil, nil
}
