package postgres

import (
	"context"
	"database/sql"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresEmailVerificationRepository struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) repository.EmailVerificationRepository {
	return &postgresEmailVerificationRepository{
		db: db,
	}
}

func (pevr *postgresEmailVerificationRepository) Save(
	ctx context.Context,
	userId uuid.UUID,
	hash string,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO email_verifications (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at;
	`

	_, err := pevr.db.ExecContext(ctx, query, userId, hash, expiresAt)

	return err
}

func (pevr *postgresEmailVerificationRepository) GetByTokenHash(
	ctx context.Context,
	hash string,
) (emailVerification *models.EmailVerification, err error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM email_verifications
		WHERE token_hash = $1
		LIMIT 1;
	`

	emailVerification = new(models.EmailVerification)
	err = pevr.db.QueryRowContext(ctx, query, hash).Scan(
		&emailVerification.ID,
		&emailVerification.UserId,
		&emailVerification.Hash,
		&emailVerification.ExpiresAt,
		&emailVerification.CreatedAt,
	)

	return
}
