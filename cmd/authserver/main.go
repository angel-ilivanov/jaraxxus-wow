package main

import "wow-server/internal/network"

func main() {
	network.Start(":3724")
}
