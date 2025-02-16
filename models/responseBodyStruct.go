package models

import "database/sql"

type ResponseBody struct {
	Result  sql.Result
	Rows    sql.Rows
	Columns []string
}
