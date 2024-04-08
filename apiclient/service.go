package apiclient

import (
	"context"
	"crypto/rand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Collection *mongo.Collection
}

func (s Service) Get(workspace string, id string) (*ApiClient, error) {
	var apiClient ApiClient
	err := s.Collection.FindOne(context.Background(), bson.M{"_id": id, "workspaceId": workspace}).Decode(&apiClient)
	if err != nil {
		return nil, err
	}

	return &apiClient, nil
}

func (s Service) GetByToken(token string) (*ApiClient, error) {
	var apiClient ApiClient
	err := s.Collection.FindOne(context.Background(), bson.M{"token": token}).Decode(&apiClient)
	if err != nil {
		return nil, err
	}

	return &apiClient, nil
}

func (s Service) List(workspace string) ([]ApiClient, error) {
	var apiClients []ApiClient
	cursor, err := s.Collection.Find(context.Background(), bson.M{"workspaceId": workspace})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &apiClients)
	if err != nil {
		return nil, err
	}

	return apiClients, nil
}

func (s Service) Create(workspace string) (string, error) {
	token, err := s.generateToken(32)
	if err != nil {
		return "", err
	}

	_, err = s.Collection.InsertOne(context.Background(), ApiClient{
		WorkspaceID: workspace,
		Token:       token,
	})
	return token, err
}

func (s Service) Delete(workspace string, id string) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.M{"_id": id, "workspaceId": workspace})
	return err
}

func (s Service) generateToken(length int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
