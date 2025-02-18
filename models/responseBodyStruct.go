package models

import (
	"database/sql"
	"log"
	"net/http"
)

type ResponseBody struct {
	RequestBody *RequestBody

	// Info for 'CREATE', 'UPDATE' and 'DELETE'
	Result       sql.Result
	LastID       int64
	RowsAffected int64

	// Info for 'READ'
	Rows        *sql.Rows
	Columns     []string
	ColumnTypes []*sql.ColumnType
	Entries     []Entry // []map[string]string

	// Error handling
	Error error // ?
}

type Entry map[string]any

func NewResponseBody(req *RequestBody) *ResponseBody {
	rb := new(ResponseBody)
	rb.RequestBody = req

	return rb
}

func (rb *ResponseBody) ProcessResult() error {
	if rb.Result != nil {
		var err error

		if rb.RequestBody.OperationType == "CREATE" || rb.RequestBody.Method == http.MethodPost {
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

func (rb *ResponseBody) ProcessRows() error {
	if rb.Rows != nil {
		var err error
		rb.Columns, err = rb.Rows.Columns()
		if err != nil {
			log.Printf("proccing *sql.Rows.Columns(): %s", err)

			return err
		}

		for rb.Rows.Next() {
			entry := make(Entry)
			values := make([]interface{}, len(rb.Columns))
			valuePtrs := make([]interface{}, len(rb.Columns))

			for i := range rb.Columns {
				valuePtrs[i] = &values[i]
			}

			if err := rb.Rows.Scan(valuePtrs...); err != nil {
				log.Printf("scanning *sql.Rows: %s", err)
				return err
			}

			for i, col := range rb.Columns {
				val := values[i]
				if val == nil {
					entry[col] = ""
				} else {
					entry[col] = string(val.([]byte))
				}
			}

			rb.Entries = append(rb.Entries, entry)
		}

		if err := rb.Rows.Err(); err != nil {
			log.Printf("processing *sql.Rows: %s", err)
			return err
		}
	}

	return nil
}

func (rb *ResponseBody) getColumnMetadata(rows *sql.Rows) error {
	var err error

	rb.Columns, err = rows.Columns()
	if err != nil {
		log.Printf("fetching metadata - *sql.Rows.Columns(): %s", err)
		return err
	}

	rb.ColumnTypes, err = rows.ColumnTypes()
	if err != nil {
		log.Printf("fetching metadata - *sql.Rows.ColumnTypes(): %s", err)
		return err
	}

	return nil
}
