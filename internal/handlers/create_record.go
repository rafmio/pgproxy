package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"pgproxy/internal/utils"
	"pgproxy/pkg/dbops"
	"strings"
)

type CreateParams struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
	Values    []string `json:"values"`
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	if isValidMethod := utils.CheckHttpMethod(w, r); !isValidMethod {
		fmt.Println("CreateRecord(): method is not valide")
		return
	}

	params, err := decodeCreateRequest(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validateCreateRecord(params); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	query, values := generateInsertSQL(params)

	db, err := dbops.Connect()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := dbops.InsertToDb(db, query, values)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to insert record: %v", err))
	}

	results, err := processResult(result)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{"message": "Record created, but unable to retrieve full details"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"message":          "Record created successfully",
		"last_inserted_id": fmt.Sprintf("%d", results[0]),
		"rows_affected":    fmt.Sprintf("%d", results[1]),
	}
	json.NewEncoder(w).Encode(resp)
}

func decodeCreateRequest(r *http.Request) (*CreateParams, error) {
	defer r.Body.Close()

	if r.Body == nil {
		err := fmt.Errorf("empty request body")
		return nil, err
	}

	var params CreateParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		err := fmt.Errorf("failed to decode request body: %w", err)
		return nil, err
	}

	return &params, nil
}

func validateCreateRecord(params *CreateParams) error {
	if params.TableName == "" {
		err := fmt.Errorf("table_name is required")
		return err
	}

	if len(params.Columns) == 0 || len(params.Values) == 0 {
		err := fmt.Errorf("columns and values are required")
		return err
	}

	if len(params.Columns) != len(params.Values) {
		err := fmt.Errorf("number of columns and values must be equal")
		return err
	}

	for i, column := range params.Columns {
		if column == "" {
			return fmt.Errorf("column at index %d is empty", i)
		}
	}

	for i, value := range params.Values {
		if value == "" {
			return fmt.Errorf("value at index %d is empty", i)
		}
	}
	return nil
}

func generateInsertSQL(params *CreateParams) (string, []any) {
	columns := strings.Join(params.Columns, ", ")
	placeholders := strings.Trim(strings.Repeat("?, ", len(params.Values)), ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", params.TableName, columns, placeholders)

	var values []any
	for _, value := range params.Values {
		values = append(values, value)
	}

	return query, values
}

func processResult(result sql.Result) ([2]int64, error) {
	lastInsrtId, err := result.LastInsertId()
	if err != nil {
		return [2]int64{}, fmt.Errorf("LastInsertId: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return [2]int64{}, fmt.Errorf("RowsAffected: %v", err)
	}

	results := [2]int64{lastInsrtId, rowsAffected}

	return results, nil
}
