package main

import "github.com/angel-ilivanov/wow-server/internal/authserver"

func main() {
	authserver.Start(":3724")
}
