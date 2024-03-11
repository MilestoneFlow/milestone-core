package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/go-chi/chi/v5"
	"milestone_core/server"
	"net/http"
)

type AuthResource struct {
}

type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Routes creates a REST router for the todos resource
func (rs AuthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Auth)

	return r
}

func (rs AuthResource) Auth(w http.ResponseWriter, r *http.Request) {
	var updateInput AuthBody
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	token, err := rs.loginCognito(updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, token)
}

func (rs AuthResource) loginCognito(body AuthBody) (*string, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider("AKIAZI2LEQP6NFJQJTXC", "8yKvLgCpqbmSUa0FDWqJDHDpN5rmBgn2qXtyWo0D", "TOKEN"),
	), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)
	clientId := "7k3a12slg179el5423qb7d8iqk"

	mac := hmac.New(sha256.New, []byte("e48eirkgrrp0lufq8fil8eooam92166tpks4fl2f6else1rkua6"))
	mac.Write([]byte(body.Email + clientId))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	user, err := client.AdminInitiateAuth(context.TODO(), &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   types.AuthFlowTypeUserPasswordAuth,
		ClientId:   aws.String(clientId),
		UserPoolId: aws.String("us-east-1_zrIqQshjP"),
		AuthParameters: map[string]string{
			"USERNAME":    body.Email,
			"PASSWORD":    body.Password,
			"SECRET_HASH": secretHash,
		},
	})
	if err != nil {
		return nil, err
	}

	return user.AuthenticationResult.AccessToken, nil
}
