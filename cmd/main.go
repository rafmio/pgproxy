package main

import (
	"log"
	"pgproxy/internal/transport"
)

func main() {
	if err := transport.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
