package queries

import (
	"fmt"
	"strings"
)

// function buildExistsQuery generates an SQL query string that is designed to query the
// PostgreSQL database to see if an entry exists in the database.
func (req *requestBody) buildExistsQuery() (string, error) {
	if len(req.Params) != len(req.Columns) {
		return "", fmt.Errorf("params and columns must have the same length")
	}

	var conditions []string

	// conditions for WHERE
	for i, column := range req.Columns {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, i+1))
	}

	query := fmt.Sprintf(
		"SELECT EXISTS (SELECT 1 FROM %s WHERE %s)",
		req.TableName,
		strings.Join(conditions, " AND "),
	)

	return query, nil
}
