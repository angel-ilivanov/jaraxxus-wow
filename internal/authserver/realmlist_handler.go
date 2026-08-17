package authserver

import (
	"context"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

func (handler *RequestHandler) handleRealmList(ctx context.Context, session *AuthSession, request protocol.RealmListRequest) (protocol.Response, error) {
	return protocol.RealmListResponse{NumChars: 0}, nil // TODO: Make db query for number of characters
}
