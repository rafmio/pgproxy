package dbops

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// import the PostgreSQL driver for database/sql
	_ "github.com/lib/pq" // $ go get .
)

type DBConfig struct {
	DriverName string
	Host       string
	Port       string
	DBName     string
	User       string
	Password   string
	SslMode    string
}

// In Go, it is customary to separate configuration and state.
// DB in this case is a state, so I've separated it into a separate structure.
type DBConnection struct {
	db *sql.DB
}

func NewDBConfig() (*DBConfig, error) {
	cfg := &DBConfig{
		DriverName: "postgres",
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		SslMode:    os.Getenv("DB_SSL_MODE"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}

func (cfg *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.SslMode,
	)
}

func EstablishConnection(cfg *DBConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.DriverName, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func Connect() (*sql.DB, error) {
	cfg, err := NewDBConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	db, err := EstablishConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	return db, nil
}
