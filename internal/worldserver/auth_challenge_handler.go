package worldserver

import (
	"crypto/rand"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
)

func handleAuthChallenge(session *session) protocol.AuthChallengeServerMessage {
	serverSeed := make([]byte, 4)
	_, _ = rand.Read(serverSeed)
	session.serverSeed = serverSeed
	return protocol.AuthChallengeServerMessage{ServerSeed: serverSeed}
}
