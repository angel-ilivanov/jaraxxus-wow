package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/database"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver"
)

func main() {
	ctx := context.Background()
	connectionString := os.Getenv("DATABASE_URL")
	if len(connectionString) == 0 {
		log.Fatal("Database connection url is missing.")
	}

	handler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

	pool, err := database.Open(ctx, connectionString)
	if err != nil {
		slog.Error("failed to open database",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
	defer pool.Close()

	accountStore := accountstore.New(pool)
	characterStore := characterstore.New(pool)
	server := worldserver.New(accountStore, characterStore)
	err = server.Start(ctx, ":8085")
	if err != nil {
		slog.Error("world server stopped",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}
