package postgres

import (
	"context"
	"database/sql"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresMfaCodeRepository struct {
	db *sql.DB
}

func NewMfaCodeRepository(db *sql.DB) repository.MfaCodeRepository {
	return &postgresMfaCodeRepository{
		db: db,
	}
}

func (pmcr *postgresMfaCodeRepository) Save(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error {
	query := `
		INSERT INTO mfa_codes (user_id, code_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			code_hash = EXCLUDED.code_hash,
			expires_at = EXCLUDED.expires_at;
	`

	_, err := pmcr.db.ExecContext(ctx, query, userId, hash, expiresAt)
	return err
}

func (pmcr *postgresMfaCodeRepository) GetByCodeHash(ctx context.Context, codeHash string) (mfaCode *models.MfaCode, err error) {
	query := `
		SELECT id, user_id, code_hash, expires_at, created_at
		FROM mfa_codes
		WHERE code_hash = $1
		LIMIT 1;
	`

	mfaCode = new(models.MfaCode)
	err = pmcr.db.QueryRowContext(ctx, query, codeHash).Scan(
		&mfaCode.ID,
		&mfaCode.UserId,
		&mfaCode.Hash,
		&mfaCode.ExpiresAt,
		&mfaCode.CreatedAt,
	)

	return
}
