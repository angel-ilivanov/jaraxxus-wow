package worldserver

type session struct {
	SessionKey []byte
	Username   string
	AccountId  string
	serverSeed []byte
	clientSeed []byte
}
