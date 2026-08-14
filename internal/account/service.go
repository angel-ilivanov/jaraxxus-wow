package account

import (
	"context"
	"fmt"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/srp6"
)

type Service struct {
	accountStore *accountstore.Store
}

func (s *Service) CreateAccount(ctx context.Context, username string, password string) (int64, error) {
	salt := srp6.GenerateSalt()
	creds := accountstore.Credentials{
		Salt:     salt,
		Verifier: srp6.CalculatePasswordVerifier(username, password, salt),
	}
	id, err := s.accountStore.Insert(ctx, username, creds)
	if err != nil {
		return 0, fmt.Errorf("create account: %w", err)
	}
	return id, nil
}

func New(accountStore *accountstore.Store) *Service {
	return &Service{accountStore: accountStore}
}
