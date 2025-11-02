package repository

import (
	"context"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) error
	VerifyEmail(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string, active bool) error
}
