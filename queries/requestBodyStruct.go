package queries

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// requestBody represents the structure of the request body.
type requestBody struct {
	tableName string   `json:"table_name"`
	columns   []string `json:"columns"`
	params    []string `json:"params"`
	newParams []string `json:"new_params"`
}

type RequestBodyErrorsMap map[*requestBody]error

// NewRequestBody parses the request body and validates it.
func NewRequestBody(w http.ResponseWriter, r *http.Request) (RequestBodyErrorsMap, error) {
	var requestBodies []*requestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBodies); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot decode request body: %s", err)
		log.Printf("Cannot decode request body: %s", err)
		return nil, err
	}

	rbem := make(RequestBodyErrorsMap) // 'rbem' stands for RequestBodyErrorMap
	badIndexSlice := make([]int, 0)

	// validate request body
	for i, rb := range requestBodies { // 'rb' stands for RequestBody
		if err := rb.validateRequestBody(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Invalid request body with idx %d: %s", i, err)
			badIndexSlice = append(badIndexSlice, i)
			rbem[rb] = err
		} else {
			rbem[rb] = nil
		}
	}

	return rbem, nil
}

// validateRequestBody validates the request body based on the provided rules.
func (req *requestBody) validateRequestBody(r *http.Request) error {
	if req.tableName == "" {
		return fmt.Errorf("table_name is required")
	}

	if len(req.params) != 0 && len(req.columns) == 0 {
		return fmt.Errorf("columns are required when params are provided")
	}

	if r.Method == http.MethodPatch && len(req.newParams) == 0 {
		return fmt.Errorf("new_params are required when method is PATCH")
	}

	return nil
}

// function buildExistsQuery generates an SQL query string that is designed to query the
// PostgreSQL database to see if an entry exists in the database.
func (req *requestBody) buildExistsQuery() (string, error) {
	if len(req.params) != len(req.columns) {
		return "", fmt.Errorf("params and columns must have the same length")
	}

	var conditions []string

	// conditions for WHERE
	for i, column := range req.columns {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, i+1))
	}

	query := fmt.Sprintf(
		"SELECT EXISTS (SELECT 1 FROM %s WHERE %s)",
		req.tableName,
		strings.Join(conditions, " AND "),
	)

	return query, nil
}
