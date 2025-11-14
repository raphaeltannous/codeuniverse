package repository

import (
	"context"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type MfaCodeRepository interface {
	Save(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error

	GetByUserId(ctx context.Context, userId uuid.UUID) (mfaCode *models.MfaCode, err error)
}
