package protocol

import (
	"fmt"
	"io"
)

type RealmListRequest struct{}

func (r RealmListRequest) isRequest() {}

func (r RealmListRequest) Opcode() Opcode {
	return CmdRealmList
}

func DecodeRealmListRequest(reader io.Reader) (RealmListRequest, error) {
	packet := make([]byte, 4)
	_, err := io.ReadFull(reader, packet)
	if err != nil {
		return RealmListRequest{}, fmt.Errorf("error reading realm list request packet")
	}
	return RealmListRequest{}, nil
}
