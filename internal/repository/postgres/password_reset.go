package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresPasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) repository.PasswordResetRepository {
	return &postgresPasswordResetRepository{
		db: db,
	}
}

func (pprr *postgresPasswordResetRepository) Save(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at;
	`

	_, err := pprr.db.ExecContext(ctx, query, userId, hash, expiresAt)

	return err
}

func (pprr *postgresPasswordResetRepository) GetByTokenHash(
	ctx context.Context,
	hash string,
) (*models.PasswordReset, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM password_resets
		WHERE token_hash = $1
		LIMIT 1;
	`

	passwordReset := new(models.PasswordReset)
	err := pprr.db.QueryRowContext(ctx, query, hash).Scan(
		&passwordReset.ID,
		&passwordReset.UserId,
		&passwordReset.Hash,
		&passwordReset.ExpiresAt,
		&passwordReset.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrPasswordResetNotFound
		}

		return nil, err
	}

	return passwordReset, nil
}
