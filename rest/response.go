package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(writer http.ResponseWriter, error error, statusCode int) {
	jsonData, err := json.Marshal(ErrorResponse{Error: error.Error()})
	if err != nil {
		log.Fatal(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_, err = writer.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}
