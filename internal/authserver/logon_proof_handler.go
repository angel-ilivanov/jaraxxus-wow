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
	//calculate session key
	//verify client proof
	//send server proof

	err := updateLogonProofSession(session, request)
	if err != nil {
		return nil, fmt.Errorf("error updating session for logon proof: %w", err)
	}

	validProof := isValidClientProof(session, request)
	if !validProof {
		response := constructLogonProofResponse(protocol.ResultIncorrectPassword, nil)
		return protocol.EncodeLogonProofResponse(response), nil
	}

	err = handler.store.SetSessionKey(ctx, session.identity.accountID, session.sessionKey)
	if err != nil {
		return nil, fmt.Errorf("error updating session key in database: %w", err)
	}

	serverProof := srp6.CalculateServerProof(
		request.ClientPublicKey,
		request.ClientProof,
		session.sessionKey)

	response := constructLogonProofResponse(protocol.ResultSuccess, serverProof)
	return protocol.EncodeLogonProofResponse(response), nil //TODO
}

func constructLogonProofResponse(result protocol.LoginResult, serverProof []byte) protocol.LogonProofResponse {
	return protocol.LogonProofResponse{Result: result, ServerProof: serverProof}
}

func updateLogonProofSession(session *AuthSession, request protocol.LogonProofRequest) error {
	fmt.Println("updating session with key")
	if session.phase != PhaseAwaitingProof {
		return fmt.Errorf("session is not in the AwaitingProof phase")
	}
	sessionKey := srp6.CalculateServerSessionKey(
		request.ClientPublicKey,
		session.srpState.serverPublicKey,
		session.srpState.verifier,
		session.srpState.serverPrivateKey)
	err := session.markAuthenticated(sessionKey)
	if err != nil {
		return fmt.Errorf("error switching session state: %w", err)
	}
	fmt.Println("session successfully updated")
	return nil
}

func isValidClientProof(session *AuthSession, request protocol.LogonProofRequest) bool {
	fmt.Println("validating client proof...")

	expectedClientProof := srp6.CalculateExpectedClientProof(
		session.identity.username,
		session.sessionKey,
		request.ClientPublicKey,
		session.srpState.serverPublicKey,
		session.identity.salt)

	fmt.Println("expected proof calculated")

	if subtle.ConstantTimeCompare(request.ClientProof, expectedClientProof) == 0 {
		fmt.Println("invalid client proof")
		return false
	}
	fmt.Println("client proof confirmed to be valid")
	return true
}
