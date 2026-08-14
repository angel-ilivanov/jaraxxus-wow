package authserver

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/srp6"
)

func (s *Server) RegisterAccount(ctx context.Context, username string, password string) (int64, error) {
	salt := srp6.GenerateSalt()
	creds := accountstore.Credentials{
		Salt:     salt,
		Verifier: srp6.CalculatePasswordVerifier(username, password, salt),
	}
	// Username and password are uppercased by client
	id, err := s.accountStore.Insert(ctx, username, creds)
	if err != nil {
		return 0, fmt.Errorf("error inserting account in database: %v", err)
	}
	return id, nil
}
