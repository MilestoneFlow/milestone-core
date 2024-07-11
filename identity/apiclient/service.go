package apiclient

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DbConnection *sqlx.DB
}

func (s Service) Get(workspace string, id string) (*ApiClient, error) {
	var apiClient ApiClient
	err := s.DbConnection.Get(&apiClient, "SELECT id, workspace_id, token FROM identity.api_client WHERE id = $1 AND workspace_id = $2", id, workspace)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &apiClient, nil
}

func (s Service) GetByToken(token string) (*ApiClient, error) {
	var apiClient ApiClient
	err := s.DbConnection.Get(&apiClient, "SELECT id, workspace_id, token FROM identity.api_client WHERE token = $1", token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &apiClient, nil
}

func (s Service) List(workspace string) ([]ApiClient, error) {
	var apiClients []ApiClient
	err := s.DbConnection.Select(&apiClients, "SELECT id, workspace_id, token FROM identity.api_client WHERE workspace_id = $1", workspace)
	if err != nil {
		return nil, err
	}

	if apiClients == nil {
		apiClients = make([]ApiClient, 0)
	}

	return apiClients, nil
}

func (s Service) Create(workspace string) (string, error) {
	token, err := s.generateToken(32)
	if err != nil {
		return "", err
	}

	_, err = s.DbConnection.Exec("INSERT INTO identity.api_client (workspace_id, token) VALUES ($1, $2)", workspace, token)
	return token, err
}

func (s Service) Delete(workspace string, id string) error {
	_, err := s.DbConnection.Exec("DELETE FROM identity.api_client WHERE id = $1 AND workspace_id = $2", id, workspace)
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
