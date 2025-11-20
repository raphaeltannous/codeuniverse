package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (pmcr *postgresMfaCodeRepository) Save(ctx context.Context, userId uuid.UUID, tokenHash, codeHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO mfa_codes (user_id, token_hash, code_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			code_hash = EXCLUDED.code_hash,
			expires_at = EXCLUDED.expires_at;
	`

	_, err := pmcr.db.ExecContext(
		ctx,
		query,
		userId,
		tokenHash,
		codeHash,
		expiresAt,
	)

	return err
}

func (pmcr *postgresMfaCodeRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.MfaCode, error) {
	query := `
		SELECT id, user_id, token_hash, code_hash, expires_at, created_at
		FROM mfa_codes
		WHERE token_hash = $1
		LIMIT 1;
	`

	mfaCode := new(models.MfaCode)
	err := pmcr.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&mfaCode.ID,
		&mfaCode.UserId,
		&mfaCode.TokenHash,
		&mfaCode.CodeHash,
		&mfaCode.ExpiresAt,
		&mfaCode.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrMfaTokenNotFound
		}
		return nil, fmt.Errorf("failed to scan mfaCode data into *model.MfaCode: %w", err)
	}

	return mfaCode, nil
}
