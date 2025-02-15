package queries

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// crudEntry represents the structure of the request body.
type crudEntry struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
	Params    []string `json:"params"`
	NewParams []string `json:"new_params"`
}

type requestBody struct {
	entries     []*crudEntry
	entryErrors map[*crudEntry]error
	badIdxs     []int
}

// newRequestBody parses the request body and validates it.
func newRequestBody(w http.ResponseWriter, r *http.Request) (*requestBody, error) {
	var parsedEntries []*crudEntry

	if err := json.NewDecoder(r.Body).Decode(&parsedEntries); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot decode request body: %s", err)
		log.Printf("Cannot decode request body: %s", err)
		return nil, err
	}

	rb := new(requestBody) // 'rb' stands for 'Request Body'
	rb.entries = make([]*crudEntry, 0)
	rb.entryErrors = make(map[*crudEntry]error)
	rb.badIdxs = make([]int, 0)

	// validate request body
	log.Println("validating data...")
	for i, parsedEntry := range parsedEntries {
		if err := parsedEntry.validateRequestBody(r); err != nil {
			log.Printf("Invalid request item with idx %d: %s", i, err)
			rb.badIdxs = append(rb.badIdxs, i)
			rb.entryErrors[parsedEntry] = err
		} else {
			rb.entries = append(rb.entries, parsedEntry)
		}
	}
	log.Println("the data validation is over")

	return rb, nil
}

// validateRequestBody validates the request body based on the provided rules.
func (req *crudEntry) validateRequestBody(r *http.Request) error {
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
