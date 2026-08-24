package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/bootstrap"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver"
)

func main() {
	ctx := context.Background()
	connectionString := os.Getenv("DATABASE_URL")
	if len(connectionString) == 0 {
		log.Fatal("Database connection url is missing.")
	}

	resources, err := bootstrap.Open(ctx, connectionString)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create resources",
			slog.Any("err", err))
		os.Exit(1)
	}
	defer resources.Close()
	server := worldserver.New(resources.Accounts)
	err = server.Start(ctx, ":8085")
	if err != nil {
		slog.Error("world server stopped",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}
