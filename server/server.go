package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerConfig struct {
	IP   string
	Port string
}

func loadConfig() ServerConfig {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Println("Port environment variable not set. Using default: 8080")
		port = "8080"
	}

	ip, ok := os.LookupEnv("IP")
	if !ok {
		log.Println("IP environment variable not set. Using default: 0.0.0.0")
		ip = "0.0.0.0"
	}

	return ServerConfig{IP: ip, Port: port}
}

func NewServer() *http.Server {
	cfg := loadConfig()
	addr := fmt.Sprintf("%s:%s", cfg.IP, cfg.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      newServeMux(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func newServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	endpoints := map[string]handlerFunc{
		"/create": createRecord,
		"/read":   readRecord,
		"/update": updateRecord,
		"/delete": deleteRecord,
		"/exists": existsRecord,
	}

	for path, handler := range endpoints {
		mux.HandleFunc(path, loggingMiddleware(recoveryMiddleware(handler)))
		log.Printf("Registered endpoint: %s", path)
	}

	return mux
}

func loggingMiddleware(next handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func recoveryMiddleware(next handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// RunServer() runs server and properly stop working
func RunServer() {
	server := NewServer()

	log.Printf("Starting server on %s...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
