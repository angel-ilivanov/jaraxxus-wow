package authserver

import (
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

// handleLogonProof updates session state and dispatches to encoder
func (handler *RequestHandler) handleLogonProof(ctx context.Context, session *AuthSession, request protocol.LogonProofRequest) ([]byte, error) {
	sessionKey := calculateSessionKey(session, request)

	validProof := isValidClientProof(session, request, sessionKey)
	if !validProof {
		session.resetForLogon()
		response := protocol.LogonProofResponse{Result: protocol.ResultIncorrectPassword, ServerProof: nil}
		return protocol.EncodeLogonProofResponse(response), nil
	}

	err := session.markAuthenticated(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("error updating session for logon proof: %w", err)
	}

	err = handler.store.SetSessionKey(ctx, session.identity.accountID, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("error updating session key in database: %w", err)
	}

	serverProof := calculateServerProof(request, sessionKey)
	response := protocol.LogonProofResponse{Result: protocol.ResultSuccess, ServerProof: serverProof}
	return protocol.EncodeLogonProofResponse(response), nil
}

func isValidClientProof(session *AuthSession, request protocol.LogonProofRequest, sessionKey []byte) bool {
	expectedClientProof := srp6.CalculateExpectedClientProof(
		session.identity.username,
		sessionKey,
		request.ClientPublicKey,
		session.srpState.serverPublicKey,
		session.identity.salt)

	if subtle.ConstantTimeCompare(request.ClientProof, expectedClientProof) == 0 {
		fmt.Println("invalid client proof")
		return false
	}
	return true
}

func calculateSessionKey(session *AuthSession, request protocol.LogonProofRequest) []byte {
	return srp6.CalculateServerSessionKey(
		request.ClientPublicKey,
		session.srpState.serverPublicKey,
		session.srpState.verifier,
		session.srpState.serverPrivateKey)
}

func calculateServerProof(request protocol.LogonProofRequest, sessionKey []byte) []byte {
	return srp6.CalculateServerProof(
		request.ClientPublicKey,
		request.ClientProof,
		sessionKey)
}
