package repository

import (
	"context"
	"errors"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type MfaCodeRepository interface {
	Save(ctx context.Context, userId uuid.UUID, tokenHash, codeHash string, expiresAt time.Time) error

	GetByTokenHash(ctx context.Context, tokenHash string) (*models.MfaCode, error)
}

var (
	ErrMfaTokenNotFound = errors.New("mfa token not found")
)
