package worldserver

type serverMessage interface {
	Encode() []byte
}
