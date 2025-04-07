package utils

import (
	"fmt"
	"net/http"
	// "pgproxy/internal/utils"
)

// Path constants for the server endpoints.
const (
	pathCreate = "/create"
	pathRead   = "/read"
	pathUpdate = "/update"
	pathDelete = "/delete"
	pathHealth = "/health"
)

// matchMethods maps server paths to their allowed HTTP methods.
var matchMethods = map[string]string{
	pathCreate: http.MethodPost,
	pathRead:   http.MethodGet,
	pathUpdate: http.MethodPut,
	pathDelete: http.MethodDelete,
	pathHealth: http.MethodGet,
}

func CheckHttpMethod(w http.ResponseWriter, r *http.Request) bool {
	if method, ok := matchMethods[r.URL.Path]; !ok {
		ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Sprintf("Path %s is not allowed", r.URL.Path))
		return false
	} else if method != r.Method {
		ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method %s not allowed for path %s, should be %s", r.Method, r.URL.Path, method))
		return false
	}

	return true
}
