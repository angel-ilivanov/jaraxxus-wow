package main

import (
	"context"
	"log"
	"os"

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

	var version string
	err = pool.QueryRow(ctx, "SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalf("verify database connection: %v", err)
	}
	log.Printf("DB Connection successful: %s", version)

	authserver.Start(":3724")
}
