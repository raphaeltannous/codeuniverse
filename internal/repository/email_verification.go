package repository

import (
	"context"
	"errors"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type EmailVerificationRepository interface {
	Save(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error

	GetByTokenHash(ctx context.Context, hash string) (*models.EmailVerification, error)
}

var (
	ErrEmailVerificationNotFound = errors.New("postgres: email verification not found")
)
