package handlers

import (
	"encoding/json"
	"net/http"
)

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{Message: "Record created successfully"}
	json.NewEncoder(w).Encode(resp)
}
