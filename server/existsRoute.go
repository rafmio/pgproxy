package server

import (
	"net/http"
)

func existsRecord(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
