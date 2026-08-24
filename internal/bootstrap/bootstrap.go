package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resources struct {
	Accounts *accountstore.Store
	pool     *pgxpool.Pool // private so it can only be closed by caller
}

// Open sets up logger, creates database pool and accountStore instance
func Open(ctx context.Context, connectionString string) (Resources, error) {
	//handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	//remove comment for debug messages
	handler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))
	pool, err := database.Open(ctx, connectionString)
	if err != nil {
		return Resources{}, fmt.Errorf("open database: %v", err)
	}
	accountStore := accountstore.New(pool)
	return Resources{
		Accounts: accountStore,
		pool:     pool,
	}, nil
}

func (r Resources) Close() {
	r.pool.Close()
}
