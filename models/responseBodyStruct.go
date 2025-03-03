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
	ResultErrors []error

	// Info for 'READ'
	Rows        *sql.Rows
	ColumnNames []string
	ColumnTypes []string
	Entries     []map[string]any
	RowsErrors  []error
}

func NewResponseBody(req *RequestBody) *ResponseBody {
	rb := new(ResponseBody)
	rb.RequestBody = req

	return rb
}

func (rb *ResponseBody) ProcessResult() error {
	if rb.Result == nil {
		return errors.New("the 'Result' field is nil")
	}

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
	return nil
}

func (rb *ResponseBody) ProcessRows() error {
	if rb.Rows == nil {
		return errors.New("the 'Rows' field is nil")
	}

	// Initialize the 'ColumnNames' and 'ColumnTypes' fields of the ResponseBody structure
	err := rb.getColumnMetadata()
	if err != nil {
		return err
	}

	for rb.Rows.Next() {
		convertor := newValueConvertor(rb)

		if err := rb.Rows.Scan(convertor.valuesPtr...); err != nil {
			log.Printf("scanning *sql.Rows: %s", err)
			return err
		}

		for i, columnName := range rb.ColumnNames {
			columnValue := convertor.values[i]
			columnType := rb.ColumnTypes[i]

			if columnValue == nil {
				convertor.entry[columnName] = string("")
			} else {
				// Depending on the column type, assign the value to the entry map
				convertor.entry[columnName], err = convertDatabaseValue(
					columnType,
					columnValue,
					columnName,
				)
			}
		}

		rb.Entries = append(rb.Entries, convertor.entry)
	}

	if err := rb.Rows.Err(); err != nil {
		log.Printf("processing *sql.Rows: %s", err)
		return err
	}

	return nil
}

func (rb *ResponseBody) getColumnMetadata() error {
	var err error
	errorHandler := func(error) error {
		err = errors.New(fmt.Sprintf("fetching DB metadata: %s", err))
		log.Println(err)
		return err
	}

	rb.ColumnNames, err = rb.Rows.Columns()
	if err != nil {
		err = errorHandler(err)
		return err
	}

	columnTypes, err := rb.Rows.ColumnTypes()
	if err != nil {
		err = errorHandler(err)
		return err
	}

	// convert []*ColumnType to []string
	for _, ct := range columnTypes {
		ctString := ct.DatabaseTypeName()
		rb.ColumnTypes = append(rb.ColumnTypes, ctString)
	}

	return nil
}
