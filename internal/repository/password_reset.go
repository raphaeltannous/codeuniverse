package repository

import (
	"context"
	"errors"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type PasswordResetRepository interface {
	Save(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error

	GetByTokenHash(ctx context.Context, hash string) (*models.PasswordReset, error)
}

var (
	ErrPasswordResetNotFound = errors.New("repository: password reset not found")
)
