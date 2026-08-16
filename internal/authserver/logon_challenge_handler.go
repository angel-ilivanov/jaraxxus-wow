package authserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

// handleLogonChallenge updates session state and dispatches to encoder
func (handler *RequestHandler) handleLogonChallenge(ctx context.Context, session *AuthSession, request protocol.LogonChallengeRequest) ([]byte, error) {
	creds, err := handler.store.FindForAuthentication(ctx, request.AccountName)
	if err != nil {
		if errors.Is(err, accountstore.ErrNotFound) {
			return protocol.EncodeLogonChallengeResponse(protocol.LogonChallengeResponse{Result: protocol.ResultUnknownAccount}), nil
		}
		return nil, fmt.Errorf("fetch authentication data: %w", err)
	}
	identity := accountIdentity{
		accountID: creds.ID,
		username:  request.AccountName,
		salt:      creds.Salt,
	}
	state := generateLogonChallengeSRPState(creds)
	err = updateLogonChallengeSession(session, state, identity)
	if err != nil {
		return nil, fmt.Errorf("update authentication session: %w", err)
	}
	response := constructLogonChallengeResponse(session.srpState.serverPublicKey, creds.Salt)
	packet := protocol.EncodeLogonChallengeResponse(response)
	return packet, nil
}

func generateLogonChallengeSRPState(creds accountstore.Authentication) *srpState {
	serverPrivateKey := srp6.GenerateServerPrivateKey()
	serverPublicKey := srp6.CalculateServerPublicKey(creds.Verifier, serverPrivateKey)
	return &srpState{
		salt:             creds.Salt,
		verifier:         creds.Verifier,
		serverPrivateKey: serverPrivateKey,
		serverPublicKey:  serverPublicKey,
	}
}
func constructLogonChallengeResponse(serverPublicKey []byte, salt []byte) protocol.LogonChallengeResponse {
	return protocol.LogonChallengeResponse{
		Result:          protocol.ResultSuccess,
		ServerPublicKey: serverPublicKey,
		Salt:            salt,
	}
}

func updateLogonChallengeSession(session *AuthSession, state *srpState, identity accountIdentity) error {
	if session.phase != PhaseAwaitingChallenge {
		return fmt.Errorf("session is not in the AwaitingChallenge phase")
	}
	err := session.beginProof(identity, state)
	if err != nil {
		return fmt.Errorf("error when switching state: %w", err)
	}
	return nil
}
