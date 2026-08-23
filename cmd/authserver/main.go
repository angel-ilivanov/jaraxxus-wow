package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/database"
)

func main() {
	ctx := context.Background()
	connectionString := os.Getenv("DATABASE_URL")
	if len(connectionString) == 0 {
		log.Fatal("Database connection url is missing.")
	}
	pool, err := database.Open(ctx, connectionString)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer pool.Close()

	// handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	// remove comment for debug messages
	handler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

	accountStore := accountstore.New(pool)
	server := authserver.New(accountStore)
	err = server.Start(ctx, ":3724")
	if err != nil {
		slog.Error("auth server stopped",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}
