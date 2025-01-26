package main

import (
	"dbproxy/server"
	"log"
)

func main() {
	server.RunServer()

	log.Println("Server is running...")
}
