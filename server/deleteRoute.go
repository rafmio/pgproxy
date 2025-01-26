package server

import (
	"fmt"
	"net/http"
)

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	// Implement record deletion logic here
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Deleted record")
}
