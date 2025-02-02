package server

import (
	"fmt"
	"log"
	"net/http"
)

func checkHttpMethod(w http.ResponseWriter, r *http.Request) error {
	matchMethods := map[string]string{
		"/create": http.MethodPost,
		"/read":   http.MethodGet,
		"/update": http.MethodPut,
		"/delete": http.MethodDelete,
		"/exists": http.MethodGet,
	}

	for path, method := range matchMethods {
		if path == r.URL.Path && method != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method %s not allowed for path %s, should be %s", r.Method, path, method)
			log.Printf("Method %s not allowed for path %s, should be %s", r.Method, path, method)
			return fmt.Errorf("method not allowed")
		}
	}

	return nil
}
