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

func SendResponse(writer http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)
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

func SendMessageResponse(writer http.ResponseWriter, message string, statusCode int) {
	jsonData, err := json.Marshal(map[string]string{"message": message})
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
