package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func SendJson(writer http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func SendMessageJson(writer http.ResponseWriter, message string) {
	jsonData, err := json.Marshal(MessageResponse{
		Message: message,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func SendBadRequestErrorJson(writer http.ResponseWriter, error error) {
	jsonData, err := json.Marshal(ErrorResponse{Error: error.Error()})
	if err != nil {
		log.Fatal(err)
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}
