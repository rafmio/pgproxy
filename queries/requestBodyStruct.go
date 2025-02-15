package queries

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// requestBody represents the structure of the request body.
type requestBody struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
	Params    []string `json:"params"`
	NewParams []string `json:"new_params"`
}

// validateRequestBody validates the request body based on the provided rules
type requestBodyValidationMap map[*requestBody]error

// NewRequestBody parses the request body and validates it.
func NewRequestBody(w http.ResponseWriter, r *http.Request) (requestBodyValidationMap, error) {
	var requestBodies []*requestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBodies); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot decode request body: %s", err)
		log.Printf("Cannot decode request body: %s", err)
		return nil, err
	}

	rbvm := make(requestBodyValidationMap) // 'rbvm' stands for RequestBodyValidationMap
	badIndexSlice := make([]int, 0)

	// validate request body
	log.Println("validating data...")
	for i, rb := range requestBodies { // 'rb' stands for RequestBody
		if err := rb.validateRequestBody(r); err != nil {
			log.Printf("Invalid request body with idx %d: %s", i, err)
			badIndexSlice = append(badIndexSlice, i)
			rbvm[rb] = err
		} else {
			rbvm[rb] = nil
		}
	}
	log.Println("the data validation is over")

	return rbvm, nil
}

// validateRequestBody validates the request body based on the provided rules.
func (req *requestBody) validateRequestBody(r *http.Request) error {
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
