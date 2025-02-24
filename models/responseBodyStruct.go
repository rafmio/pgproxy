package models

import (
	"database/sql"
	"errors"
	"fmt"
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
		// we receive information about column names and types of received data (DB data types)
		err := rb.getColumnMetadata()
		if err != nil {
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
				columnType := rb.ColumnTypes[i].DatabaseTypeName()

				if val == nil {
					// entry[col] = emptyValues[columnType]
					fillNilValues(entry, col, columnType)
					// continue // rb.Entries slice appending will not be done?
				} else {
					// Depending on the column type, assign the value to the entry map
					entry[col], err = convertDatabaseValue(columnType, val, col)
					// switch columnType {
					// case "INT", "INTEGER", "SMALLINT", "BIGINT", "SERIAL", "BIGSERIAL":
					// 	entry[col] = val.(int64)
					// case "FLOAT", "DOUBLE PRECISION", "REAL", "NUMERIC", "DECIMAL":
					// 	entry[col] = val.(float64)
					// case "BOOLEAN":
					// 	entry[col] = val.(bool)
					// case "VARCHAR", "TEXT", "CHAR":
					// 	entry[col] = string(val.([]byte))
					// case "DATE", "TIME", "TIMESTAMP", "TIMESTAMPTZ":
					// 	entry[col] = val.(time.Time).Format(time.RFC3339)
					// case "INET", "CIDR":
					// 	entry[col] = net.IP(val.([]byte)).String()
					// default:
					// 	entry[col] = string(val.([]byte)) // Fallback to string for unknown types
					// }
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

func (rb *ResponseBody) getColumnMetadata() error {
	var err error

	rb.Columns, err = rb.Rows.Columns()
	if err != nil {
		err = errors.New(fmt.Sprintf("fetching metadata - *sql.Rows.Columns() %s", err))
		log.Println(err)
		return err
	}

	rb.ColumnTypes, err = rb.Rows.ColumnTypes()
	if err != nil {
		err = errors.New(fmt.Sprintf("fetching metadata - *sql.Rows.ColumnTypes() %s", err))
		log.Println(err)
		return err
	}

	return nil
}
