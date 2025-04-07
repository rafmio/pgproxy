package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pgproxy/internal/utils"
)

type UpdateParams struct {
	TableName string            `json:"table_name"`
	Columns   []string          `json:"columns"`
	NewValues []string          `json:"new_values"`
	Filters   map[string]string `json:"filters"`
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	isValidMethod := utils.CheckHttpMethod(w, r)
	if !isValidMethod {
		fmt.Println("CreateRecord(): method is not valide")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := Response{Message: "Record updated successfully"}
	json.NewEncoder(w).Encode(resp)
}
