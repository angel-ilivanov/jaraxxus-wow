package authserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

// handleLogonChallengeMessage performs the logon-challenge use case
func (handler *MessageHandler) handleLogonChallengeMessage(ctx context.Context, session *AuthSession, message protocol.LogonChallengeClientMessage) error {
	// validate session phase
	// query store
	// generate srp values
	// update session
	// construct response
	if session.phase != PhaseAwaitingChallenge {
		return fmt.Errorf("session is not in the AwaitingChallenge phase")
	}
	creds, err := handler.store.FindForAuthentication(ctx, message.AccountName)
	if err != nil {
		fmt.Errorf("error fetching data for username %s: %w", message.AccountName, err)
	}
	state := generateSRPState(creds)
	identity := accountIdentity{
		accountID: creds.ID,
		username:  message.AccountName,
	}
	err = session.beginProof(identity, state)
	if err != nil {
		return fmt.Errorf("error updating session: %w", err)
	}
	response := protocol.LogonChallengeServerMessage{
		Success:         true,
		ServerPublicKey: session.temporarySrpState.serverPublicKey,
		Salt:            creds.Salt,
	}
	protocol.AssembleServerChallengePacket(response)
	return nil
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
