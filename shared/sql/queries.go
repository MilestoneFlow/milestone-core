package sql

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

func FetchOne[T any](dbConnection *sqlx.DB, query string, args ...interface{}) (*T, error) {
	var result T
	err := dbConnection.Get(&result, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &result, err
}

func FetchMultiple[T any](dbConnection *sqlx.DB, query string, args ...interface{}) ([]T, error) {
	rows := make([]T, 0)
	err := dbConnection.Select(&rows, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return rows, nil
	}
	return rows, err
}
