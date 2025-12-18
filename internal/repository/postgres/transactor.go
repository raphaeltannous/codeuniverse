package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/repository"
)

const (
	txKey = "dbTx"
)

type postgresDbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type postgreSQLTransactor struct {
	db *sql.DB
}

func NewPostgreSQLTransactor(db *sql.DB) repository.Transactor {
	return &postgreSQLTransactor{db: db}
}

func (p *postgreSQLTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey, tx)

	err = fn(txCtx)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback: %w (original error: %w)", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
