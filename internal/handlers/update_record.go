package handlers

import (
	"encoding/json"
	"net/http"
)

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{Message: "Record updated successfully"}
	json.NewEncoder(w).Encode(resp)
}
