package handlers

import (
	"encoding/json"
	"net/http"
)

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{Message: "Record deleted successfully"}
	json.NewEncoder(w).Encode(resp)
}
