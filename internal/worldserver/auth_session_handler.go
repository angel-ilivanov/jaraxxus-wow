package worldserver

import (
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func (handler RequestHandler) handleAuthSession(ctx context.Context, session *session, request protocol.AuthSessionRequest) (protocol.AuthResponse, error) {
	sessionKey, err := handler.store.FetchSessionKey(ctx, request.Username)
	if err != nil {
		return protocol.AuthResponse{}, fmt.Errorf("error fetching sessionKey: %w", err)
	}
	err = updateSession(session, request, sessionKey)
	if err != nil {
		return protocol.AuthResponse{}, fmt.Errorf("error updating session: %w", err)
	}

	if !isValidClientProof(session, request, sessionKey) {
		return protocol.AuthResponse{ResultCode: protocol.ResultAuthReject}, nil
	}

	return protocol.AuthResponse{ResultCode: protocol.ResultAuthOk}, nil
}

func updateSession(session *session, request protocol.AuthSessionRequest, sessionKey []byte) error {
	session.clientSeed = request.ClientSeed
	session.Username = request.Username

	session.SessionKey = sessionKey
	return nil
}

func isValidClientProof(session *session, request protocol.AuthSessionRequest, sessionKey []byte) bool {
	expectedProof := srp6.CalculateWorldServerProof(
		request.Username,
		request.ClientSeed,
		session.serverSeed,
		sessionKey)

	if subtle.ConstantTimeCompare(request.ClientProof, expectedProof) == 0 {
		return false
	}
	return true
}
