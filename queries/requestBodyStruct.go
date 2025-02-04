package queries

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestBody struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
	Params    []string `json:"params"`
	NewParams []string `json:"new_params"`
}

func NewRequestBody(w http.ResponseWriter, r *http.Request) ([]*RequestBody, error) {
	var requestBody []*RequestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot decode request body: %s", err)
		log.Printf("Cannot decode request body: %s", err)
		return nil, err
	}

	if err := requestBody.validateRequestBody(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %s", err)
		log.Printf("Invalid request body: %s", err)
		return nil, err
	}

	return requestBody, nil
}

func (req *RequestBody) validateRequestBody(r *http.Request) error {
	if req.TableName == "" {
		return fmt.Errorf("table_name is required")
	}

	if len(req.Params) != 0 && len(req.Columns) == 0 {
		return fmt.Errorf("columns are required when params are provided")
	}

	if r.Method == http.MethodPatch && len(req.NewParams) == 0 {
		return fmt.Errorf("new_params are required when method is PATCH")
	}

	return nil
}

// function buildExistsQuery generates an SQL query string that is designed to query the
// PostgreSQL database to see if an entry exists in the database.
func (req *RequestBody) buildExistsQuery() (string, error) {
	if len(req.Params) == 0 {
		return "", fmt.Errorf("params are required to build exists query")
	}

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s = $%d)",
		req.TableName, req.Columns[0], len(req.Params)+1)

	return query, nil
}
