package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CrudEntry represents a single CRUD operation request.
type CrudEntry struct {
	OperationType string   `json:"operation_type"` // "CREATE", "READ", "UPDATE", "DELETE"
	TableName     string   `json:"table_name"`
	Columns       []string `json:"columns,omitempty"`
	Params        []string `json:"params,omitempty"`
	NewParams     []string `json:"new_params,omitempty"`
}

// RequestBody represents the entire HTTP request body for CRUD operations.
type RequestBody struct {
	Method  string       `json:"method"`  // HTTP method (GET, POST, etc.)
	Entries []*CrudEntry `json:"entries"` // List of CRUD operations
}

// NewRequestBody parses the HTTP request body and creates a RequestBody instance.
func NewRequestBody(w http.ResponseWriter, r *http.Request) (*RequestBody, error) {
	var entries []*CrudEntry
	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot decode request body: %s", err)
		log.Printf("Cannot decode request body: %s", err)
		return nil, fmt.Errorf("cannot decode request body: %w", err)
	}

	rb := &RequestBody{
		Method:  r.Method,
		Entries: entries,
	}

	if err := rb.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Validation failed: %s", err)
		log.Printf("Validation failed: %s", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return rb, nil
}

// Validate checks if the RequestBody is valid based on the operation type.
func (rb *RequestBody) Validate() error {
	for _, entry := range rb.Entries {
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("invalid entry: %w", err)
		}
	}
	return nil
}

// Validate checks if a CrudEntry is valid based on its operation type.
func (entry *CrudEntry) Validate() error {
	// Validate TableName
	if entry.TableName == "" {
		return fmt.Errorf("table_name is required")
	}

	// Validate based on OperationType
	switch entry.OperationType {
	case "CREATE":
		if len(entry.Columns) == 0 || len(entry.Params) == 0 {
			return fmt.Errorf("columns and params are required for CREATE")
		}
	case "READ":
		if len(entry.Columns) == 0 {
			return fmt.Errorf("columns are required for READ")
		}
	case "UPDATE":
		if len(entry.NewParams) == 0 {
			return fmt.Errorf("new_params are required for UPDATE")
		}
	case "DELETE":
		if len(entry.Params) == 0 {
			return fmt.Errorf("params are required for DELETE")
		}
	default:
		return fmt.Errorf("invalid operation type: %s", entry.OperationType)
	}

	return nil
}
