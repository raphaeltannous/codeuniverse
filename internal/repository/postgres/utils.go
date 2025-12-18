package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func updateColumnValue(
	ctx context.Context,
	db *sql.DB,
	table string,
	id uuid.UUID,
	column string,
	value any,
) error {
	query := fmt.Sprintf(
		`
			UPDATE %s
			set %s = $1
			WHERE id = $2;
		`,
		table,
		column,
	)

	if s, ok := value.(string); ok && s == "" {
		value = nil
	}

	result, err := db.ExecContext(
		ctx,
		query,
		value,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update %s.%s: %w", table, column, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows from database: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("expect single row affected, got %d rows affected", rows)
	}

	return nil
}

func getExecutor(ctx context.Context, db *sql.DB) postgresDbExecutor {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}

	return db
}
