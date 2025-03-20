package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

type config struct {
	ip           string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
	maxBodyBytes int64
}

const (
	defaultIP           = "0.0.0.0"
	defaultPort         = "8080"
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 120 * time.Second
	defaultMaxBodyBytes = 1_048_576 // 1MB
)

func (c *config) addr() string {
	return net.JoinHostPort(c.ip, c.port)
}

func (c *config) validate() error {
	if _, err := net.ResolveTCPAddr("tcp", c.addr()); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	return nil
}

func loadConfig() (*config, error) {
	cfg := &config{
		ip:           getEnv("IP", defaultIP),
		port:         getEnv("PORT", defaultPort),
		readTimeout:  parseDurationEnv("READ_TIMEOUT", defaultReadTimeout),
		writeTimeout: parseDurationEnv("WRITE_TIMEOUT", defaultWriteTimeout),
		idleTimeout:  parseDurationEnv("IDLE_TIMEOUT", defaultIdleTimeout),
		maxBodyBytes: parseInt64Env("MAX_BODY_BYTES", defaultMaxBodyBytes),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Invalid duration format for %s: %v. Using default", key, err)
			return defaultValue
		}
		return duration
	}
	return defaultValue
}

func parseInt64Env(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Invalid integer format for %s: %v. Using default", key, err)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

type middleware func(http.Handler) http.Handler

func chainMiddleware(h http.Handler, m ...middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Printf("PANIC: %v\n%s", err, stack)
				errorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func newServeMux(cfg *config) *http.ServeMux {
	mux := http.NewServeMux()
	endpoints := map[string]http.HandlerFunc{
		"/create": createRecord,
		"/read":   readRecord,
		"/update": updateRecord,
		"/delete": deleteRecord,
		"/health": healthCheck,
	}

	for path, handler := range endpoints {
		// http.MaxBytesHandler wraps an HTTP handler to enforce a maximum request body size.
		// If the request body exceeds the specified limit, the handler will return an error
		// with a 413 status code (Request Entity Too Large), protecting the server from
		// large payloads that could lead to resource exhaustion or denial-of-service attacks.
		wrappedHandler := http.MaxBytesHandler(handler, cfg.maxBodyBytes)

		mux.Handle(
			path,
			chainMiddleware(
				wrappedHandler,
				loggingMiddleware,
				recoveryMiddleware,
			),
		)
		log.Printf("Registered endpoint: %s", path)
	}
	return mux
}

func newServer(cfg *config) *http.Server {
	return &http.Server{
		Addr:         cfg.addr(),
		Handler:      newServeMux(cfg),
		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,
		IdleTimeout:  cfg.idleTimeout,
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Run запускает HTTP сервер с graceful shutdown
func Run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	server := newServer(cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil

}

// Добавленные недостающие обработчики
func createRecord(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать логику создания записи
	errorResponse(w, http.StatusNotImplemented, "Not implemented")
}

func readRecord(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать логику чтения записи
	errorResponse(w, http.StatusNotImplemented, "Not implemented")
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать логику обновления записи
	errorResponse(w, http.StatusNotImplemented, "Not implemented")
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать логику удаления записи
	errorResponse(w, http.StatusNotImplemented, "Not implemented")
}
