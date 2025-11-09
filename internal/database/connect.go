package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrEmptyDBString = errors.New("db string is empty")
)

func Connect() (*sql.DB, error) {
	connectionString := os.Getenv("CODEUNIVERSE_DBSTRING")
	if connectionString == "" {
		return nil, ErrEmptyDBString
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to make a connection to the datbase: %w", err)
	}

	return db, nil
}
