package dbops

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
	Dsn        string
	DB         *sql.DB
}

func NewDBConfig() *DBConfig {
	dbCfg := new(DBConfig)
	dbCfg.DriverName = "postgres"
	dbCfg.Host = os.Getenv("DB_HOST")
	dbCfg.Port = os.Getenv("DB_PORT")
	dbCfg.User = os.Getenv("DB_USER")
	dbCfg.Password = os.Getenv("DB_PASSWORD")
	dbCfg.DBName = os.Getenv("DB_NAME")
	dbCfg.SslMode = os.Getenv("DB_SSL_MODE")

	log.Println("new DB config has been created")

	return dbCfg
}

func (dbc *DBConfig) SetDSN() {
	formatString := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"

	dbc.Dsn = fmt.Sprintf(formatString,
		dbc.Host,
		dbc.Port,
		dbc.User,
		dbc.DBName,
		dbc.Password,
		dbc.SslMode,
	)
	log.Println("DSN has been set")
}

func (dbc *DBConfig) EstablishDbConnection() error {
	var err error
	dbc.DB, err = sql.Open(dbc.DriverName, dbc.Dsn)
	if err != nil {
		log.Println("Open database:", err)
		return err
	}

	err = dbc.DB.Ping()
	if err != nil {
		log.Println("Ping database:", err)
	}
	log.Println("connection to DB has been established")

	return nil
}

func ConnectToDb() (*sql.DB, error) {
	dbCfg := NewDBConfig()
	dbCfg.SetDSN()

	log.Println("Establishing connection to DB...")
	err := dbCfg.EstablishDbConnection()
	if err != nil {
		return nil, err
	}

	return dbCfg.DB, nil
}
