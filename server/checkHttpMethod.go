package server

import (
	"fmt"
	"net/http"
)

// Path constants for the server endpoints.
const (
	pathCreate = "/create"
	pathRead   = "/read"
	pathUpdate = "/update"
	pathDelete = "/delete"
	pathExists = "/exists"
)

// matchMethods maps server paths to their allowed HTTP methods.
var matchMethods = map[string]string{
	pathCreate: http.MethodPost,
	pathRead:   http.MethodGet,
	pathUpdate: http.MethodPut,
	pathDelete: http.MethodDelete,
	pathExists: http.MethodGet,
}

func checkHttpMethod(w http.ResponseWriter, r *http.Request) bool {
	sendError := func(statusCode int, format string, args ...interface{}) bool {
		http.Error(w, fmt.Sprintf(format, args...), statusCode)
		return false
	}

	expectedMethod, ok := matchMethods[r.URL.Path]
	if !ok {
		return sendError(http.StatusMethodNotAllowed, "Path %s is not allowed", r.URL.Path)
	}

	if expectedMethod != r.Method {
		return sendError(http.StatusMethodNotAllowed, "Method %s not allowed for path %s, should be %s",
			r.Method, r.URL.Path, expectedMethod)
	}

	return true
}
