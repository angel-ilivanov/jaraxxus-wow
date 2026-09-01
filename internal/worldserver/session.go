package worldserver

type session struct {
	SessionKey  []byte
	AccountName string
	AccountId   string
	serverSeed  []byte
	clientSeed  []byte
}
