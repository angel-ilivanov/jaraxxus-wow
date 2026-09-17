package protocol

type TutorialFlagsResponse struct {
}

func (t TutorialFlagsResponse) EncodeBody() []byte {
	return make([]byte, 32) // all tutorials not passed
}

func (t TutorialFlagsResponse) Opcode() ServerOpcode {
	return ServerOpcodeTutorialFlags
}
