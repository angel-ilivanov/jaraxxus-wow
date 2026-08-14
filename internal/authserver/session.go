package authserver

import "fmt"

type AuthSession struct {
	phase             authPhase
	identity          *accountIdentity
	temporarySrpState *srpState
	sessionKey        []byte
}

type authPhase uint8 //enum
const (
	PhaseAwaitingChallenge authPhase = iota
	PhaseAwaitingProof
	PhaseAuthenticated
)

type accountIdentity struct {
	accountID int64
	username  string
}
type srpState struct {
	salt             []byte
	verifier         []byte
	serverPrivateKey []byte
	serverPublicKey  []byte
}

func (s *AuthSession) beginProof(identity *accountIdentity, state *srpState) error {
	if s.phase != PhaseAwaitingChallenge {
		return fmt.Errorf("session proof has already started")
	}
	s.identity = identity
	s.temporarySrpState = state
	s.phase = PhaseAwaitingProof
	return nil
}

func (s *AuthSession) markAuthenticated(sessionKey []byte) error {
	if s.phase != PhaseAwaitingProof {
		return fmt.Errorf("session cannot be authenticated at phase %v", s.phase)
	}
	s.sessionKey = sessionKey
	s.temporarySrpState = nil
	s.phase = PhaseAuthenticated
	return nil
}
