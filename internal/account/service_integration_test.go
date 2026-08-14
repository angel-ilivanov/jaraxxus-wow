package account

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/database"
)

func TestRegisterAccount_Success(t *testing.T) {
	ctx := context.Background()
	service := setup(ctx, t)

	username := "acc_" + rand.Text()[0:6]
	password := "MVRVMUJFWRA0IBVK"

	id, err := service.CreateAccount(ctx, username, password)
	if err != nil {
		t.Fatalf("error registering account: %v", err)
	}
	auth, err := service.accountStore.FindForAuthentication(ctx, username)
	if err != nil {
		t.Fatalf("error querying database: %v", err)
	}
	if id != auth.ID {
		t.Fatalf("accountId = %v, expected %v", auth.ID, id)
	}
}
func TestRegisterAccount_UsernameTaken(t *testing.T) {
	ctx := context.Background()
	service := setup(ctx, t)

	username := "acc_" + rand.Text()[0:6]
	password := "MVRVMUJFWRA0IBVK"

	_, err := service.CreateAccount(ctx, username, password)
	if err != nil {
		t.Fatalf("error registering account: %v", err)
	}
	_, err = service.CreateAccount(ctx, username, password)
	if err == nil {
		t.Fatalf("expected an error when registering with a taken username")
	}
}

func setup(ctx context.Context, t *testing.T) *Service {
	connectionString := os.Getenv("TEST_DATABASE_URL")
	if connectionString == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}
	pool, err := database.Open(ctx, connectionString)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(pool.Close)

	store := accountstore.New(pool)
	return New(store)
}
