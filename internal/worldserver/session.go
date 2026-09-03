package worldserver

type session struct {
	SessionKey []byte
	Username   string
	AccountId  int64
	serverSeed []byte
	clientSeed []byte
}
