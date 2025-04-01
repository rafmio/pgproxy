package dbops

import (
	"database/sql"
	"log"
	"os"
)

func RunCreateTableScript(db *sql.DB) error {
	filepath := os.Getenv("DB_CREATE_TABLE_SCRIPT")
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Println("error reading create table script file", err)
		return err
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Println("error executing create table script", err)
		return err
	}

	return nil
}
