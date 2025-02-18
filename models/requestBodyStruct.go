package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CrudEntry represents the structure of the request body.
type CrudEntry struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
	Params    []string `json:"params"`
	NewParams []string `json:"new_params"`
}

type RequestBody struct {
	Method        string // http.Request.Method
	OperationType string // one of the CRUD
	Entries       []*CrudEntry
	EntryErrors   map[*CrudEntry]error
	BadIdxs       []int
}

// newRequestBody parses the request body and validates it.
func NewRequestBody(w http.ResponseWriter, r *http.Request) (*RequestBody, error) {
	var parsedEntries []*CrudEntry

	if err := json.NewDecoder(r.Body).Decode(&parsedEntries); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot decode request body: %s", err)
		log.Printf("Cannot decode request body: %s", err)
		return nil, err
	}

	rb := new(RequestBody) // 'rb' stands for 'Request Body'
	rb.Entries = make([]*CrudEntry, 0)
	rb.EntryErrors = make(map[*CrudEntry]error)
	rb.BadIdxs = make([]int, 0)

	// validate request body
	log.Println("validating data...")
	for i, parsedEntry := range parsedEntries {
		if err := parsedEntry.validateRequestBody(r); err != nil {
			log.Printf("Invalid request item with idx %d: %s", i, err)
			rb.BadIdxs = append(rb.BadIdxs, i)
			rb.EntryErrors[parsedEntry] = err
		} else {
			rb.Entries = append(rb.Entries, parsedEntry)
		}
	}
	log.Println("the data validation is over")

	return rb, nil
}

// validateRequestBody validates the request body based on the provided rules.
func (req *CrudEntry) validateRequestBody(r *http.Request) error {
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
