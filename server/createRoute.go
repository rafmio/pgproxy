package server

import (
	"dbproxy/dbops"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func createRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		log.Println("Method not allowed")
		return
	}

	var requestBody RequestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		log.Println("Invalid request body")
		return
	}

	if requestBody.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Query cannot be empty")
		log.Println("Query cannot be empty")
		return
	}

	db, err := dbops.ConnectToDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to connect to database")
		log.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	params := make([]any, len(requestBody.Params))
	for i, param := range requestBody.Params {
		params[i] = param
	}

	_, err = dbops.MakeQuery(db, requestBody.Query, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to create record")
		log.Println("Failed to create record:", err)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Record created successfully")
		log.Println("Record created successfully")
		w.WriteHeader(http.StatusOK)
	}
}
