package server

import (
	"fmt"
	"log"
	"net/http"
	"pgproxy/dbops"
	"pgproxy/models"
)

func createRecord(w http.ResponseWriter, r *http.Request) {
	err := checkHttpMethod(w, r)
	if err != nil {
		return
	}

	requestBody, err := models.NewRequestBody(w, r)
	if err != nil {
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
		// w.WriteHeader(http.StatusOK) // deleting this fix the error "superfluous response.WriteHeader call"
	}
}
