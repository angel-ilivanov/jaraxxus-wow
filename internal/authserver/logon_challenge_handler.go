package authserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

// handleLogonChallengeMessage performs the logon-challenge use case
func (handler *MessageHandler) handleLogonChallengeMessage(ctx context.Context, session *AuthSession, message protocol.LogonChallengeClientMessage) ([]byte, error) {
	creds, err := handler.store.FindForAuthentication(ctx, message.AccountName)
	if err != nil {
		return nil, fmt.Errorf("error fetching data for username %s: %w", message.AccountName, err)
	}
	identity := accountIdentity{
		accountID: creds.ID,
		username:  message.AccountName,
	}
	state := generateSRPState(creds)
	err = updateSession(session, state, identity)
	if err != nil {
		return nil, fmt.Errorf("error updating session: %w", err)
	}
	response := constructMessage(session.temporarySrpState.serverPublicKey, creds.Salt)
	packet := protocol.AssembleServerChallengePacket(response)
	return packet, nil
}

func generateSRPState(creds accountstore.Authentication) *srpState {
	serverPrivateKey := srp6.GenerateServerPrivateKey()
	serverPublicKey := srp6.CalculateServerPublicKey(creds.Verifier, serverPrivateKey)
	return &srpState{
		salt:             creds.Salt,
		verifier:         creds.Verifier,
		serverPrivateKey: serverPrivateKey,
		serverPublicKey:  serverPublicKey,
	}
}
func constructMessage(serverPublicKey []byte, salt []byte) protocol.LogonChallengeServerMessage {
	return protocol.LogonChallengeServerMessage{
		Success:         true,
		ServerPublicKey: serverPublicKey,
		Salt:            salt,
	}
}

func updateSession(session *AuthSession, state *srpState, identity accountIdentity) error {
	if session.phase != PhaseAwaitingChallenge {
		return fmt.Errorf("session is not in the AwaitingChallenge phase")
	}
	err := session.beginProof(identity, state)
	if err != nil {
		return fmt.Errorf("error when switching state: %w", err)
	}
	return nil
}
