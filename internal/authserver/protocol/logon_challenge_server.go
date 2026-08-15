package protocol

//define response packet
//encode it
// wire representation of the server response

type LogonChallengeServerMessage struct {
	Success         bool
	ServerPublicKey []byte
	Salt            []byte
}

func AssembleServerChallengePacket(message LogonChallengeServerMessage) {
	//TODO
}
