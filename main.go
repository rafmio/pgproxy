package main

import (
	"log"
	"pgproxy/server"
)

func main() {
	server.RunServer()

	log.Println("Server is running...")
}
