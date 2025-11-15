package repository

import (
	"context"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type EmailVerificationRepository interface {
	Save(ctx context.Context, userId uuid.UUID, hash string, expiresAt time.Time) error

	GetByTokenHash(ctx context.Context, hash string) (emailVerification *models.EmailVerification, err error)
}
