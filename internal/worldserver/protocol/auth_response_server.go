package protocol

type AuthResponse struct {
	ResultCode AccountResultValue
}

func (a AuthResponse) Encode() []byte {
	//TODO implement me
	panic("implement me")
}
