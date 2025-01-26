package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Columns []string `json:"columns"`
type RequestBody struct {
	Query  string   `json:"query"`
	Params []string `json:"params"`
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func NewServer() *http.Server {
	server := new(http.Server)

	// Set up environment variables for IP and PORT if not provided
	port := os.Getenv("PORT")
	ip := os.Getenv("IP")
	if ip == "" || port == "" {
		// Default to 8080 and 0.0.0.0 if not provided in environment variables
		log.Println("Environment variables IP and PORT not provided. Defaulting to 0.0.0.0:8080.")
		port = "8080"
		ip = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%s", ip, port)

	server.Addr = addr
	server.Handler = newServeMux() // Creating a new router and attach it to the server
	server.ReadTimeout = 10 * time.Second

	return server
}

// newServeMux() creates a new router and attaches handlers to it
func newServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	endpoints := map[string]handlerFunc{
		"/create": createRecord,
		"/read":   readRecord,
		"/delete": deleteRecord,
		"/update": updateRecord,
		"/exists": existsRecord,
	}

	for path, handler := range endpoints {
		mux.HandleFunc(path, handler)
	}

	return mux
}

// RunServer() runs server and properly stop working
func RunServer() {
	server := NewServer()

	log.Printf("Starting server on %s...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
