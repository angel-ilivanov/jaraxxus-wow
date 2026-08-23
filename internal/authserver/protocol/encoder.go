package protocol

import (
	"errors"
	"fmt"
)

var ErrUnknownResponseType = errors.New("unknown response type")

func EncodeResponse(response Response) ([]byte, error) {
	switch response.(type) {
	case LogonChallengeResponse:
		return EncodeLogonChallengeResponse(response.(LogonChallengeResponse)), nil
	case LogonProofResponse:
		return EncodeLogonProofResponse(response.(LogonProofResponse)), nil
	case RealmListResponse:
		return EncodeRealmListResponse(response.(RealmListResponse)), nil
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnknownResponseType, response)
	}
}
