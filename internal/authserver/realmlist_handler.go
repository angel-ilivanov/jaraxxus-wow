package authserver

import (
	"context"
	"fmt"
	"os"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

func (handler *RequestHandler) handleRealmList(ctx context.Context, session *AuthSession, request protocol.RealmListRequest) (protocol.Response, error) {
	worldServerAddress := os.Getenv("WORLD_SERVER_ADDRESS")
	if worldServerAddress == "" {
		fmt.Println("address not set, using development default")
		worldServerAddress = "127.0.0.1:8085" // development default
	}
	return protocol.RealmListResponse{NumChars: 0, WorldServerAddress: worldServerAddress}, nil // TODO: Make db query for number of characters
}
