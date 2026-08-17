package protocol

import "fmt"

func EncodeResponse(response Response) ([]byte, error) {
	switch response.(type) {
	case LogonChallengeResponse:
		return EncodeLogonChallengeResponse(response.(LogonChallengeResponse)), nil
	case LogonProofResponse:
		return EncodeLogonProofResponse(response.(LogonProofResponse)), nil
	default:
		return nil, fmt.Errorf("unknown response type")
	}
}
