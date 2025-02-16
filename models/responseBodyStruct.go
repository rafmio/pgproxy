package models

import (
	"database/sql"
	"net/http"
)

type ResponseBody struct {
	// Meta info
	Method        string // http.Request.Method
	OperationType string // one of the CRUD

	// Errors
	EntryErrors map[*CrudEntry]error
	BadIdxs     []int

	// Info for 'CREATE', 'UPDATE' and 'DELETE'
	Result       sql.Result
	LastID       int64
	RowsAffected int64

	// Info for 'READ'
	Rows    *sql.Rows
	Columns []string
	Entries []Entry // []map[string]string

	// Error handling
	Error error
}

type Entry map[string]string

func (rb *ResponseBody) ProcessResult() error {
	if rb.Result != nil {
		var err error

		if rb.OperationType == "CREATE" || rb.Method == http.MethodPost {
			rb.LastID, err = rb.Result.LastInsertId()
			if err != nil {
				return err
			}
		}

		rb.RowsAffected, err = rb.Result.RowsAffected()
		if err != nil {
			return err
		}
	}
	return nil
}
