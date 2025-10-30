package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	connectionString := os.Getenv("CODEUNIVERSE_DBSTRING")
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to make a connection to the datbase: %w", err)
	}

	return db, nil
}
