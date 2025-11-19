package repository

import (
	"context"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type MfaCodeRepository interface {
	Save(ctx context.Context, userId uuid.UUID, tokenHash, codeHash string, expiresAt time.Time) error

	GetByTokenHash(ctx context.Context, codeHash string) (mfaCode *models.MfaCode, err error)
}
