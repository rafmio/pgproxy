package handlers

import (
	"encoding/json"
	"net/http"
)

func ReadRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{Message: "Record read successfully"}
	json.NewEncoder(w).Encode(resp)
}
