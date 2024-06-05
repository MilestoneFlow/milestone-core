package server

import (
	"context"
	"encoding/json"
	"log"
	"milestone_core/identity/authorization"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func GetWorkspaceIdFromContext(ctx context.Context) string {
	userData := ctx.Value("user").(authorization.UserData)

	return userData.WorkspaceID
}

func GetTokenFromPublicApiClientContext(ctx context.Context) string {
	publicApiClientData := ctx.Value("client").(authorization.PublicApiClientData)

	return publicApiClientData.Token
}

func GetWorkspaceIdFromPublicApiClientContext(ctx context.Context) string {
	publicApiClientData := ctx.Value("client").(authorization.PublicApiClientData)

	return publicApiClientData.WorkspaceID
}

func SendJson(writer http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}

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

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	_, err = writer.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}
