package worldserver

type serverMessage interface {
	Encode() []byte
}

type ClientMessage interface {
	isClientMessage()
}
