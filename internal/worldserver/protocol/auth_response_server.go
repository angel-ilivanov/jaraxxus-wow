package protocol

type AuthResponse struct {
	resultCode AccountResultValue
}

func (a AuthResponse) Encode() []byte {
	//TODO implement me
	panic("implement me")
}
