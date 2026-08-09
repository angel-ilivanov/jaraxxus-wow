package main

import (
	"context"
	"log"
	"os"

	"github.com/angel-ilivanov/wow-server/internal/accountstore"
	"github.com/angel-ilivanov/wow-server/internal/authserver"
	"github.com/angel-ilivanov/wow-server/internal/database"
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

	accountStore := accountstore.New(pool)
	server := authserver.New(accountStore)
	server.Start(":3724")
}
