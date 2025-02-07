package server

import (
	"fmt"
	"log"
	"net/http"
	"pgproxy/dbops"
	"pgproxy/queries"
)

type RequestBodyErrorsMap map[queries.RequestBody]error

func createRecord(w http.ResponseWriter, r *http.Request) {
	err := checkHttpMethod(w, r)
	if err != nil {
		return
	}

	requestBody, err := queries.NewRequestBody(w, r)
	if err != nil {
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
