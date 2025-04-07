package dbops

import (
	"database/sql"
	"log"
)

// InsertToDb() inserts single entry (row), not multiple rows
func InsertToDb(db *sql.DB, query string, params []any) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("error prepare query:", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(params...)
	if err != nil {
		log.Println("error executing query:", err)
		return nil, err
	}

	return result, nil
}
