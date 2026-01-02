package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetUsers(ctx context.Context, offset, limit int) ([]*models.User, error)

	GetAdminCount(ctx context.Context) (int, error)
	GetUsersCount(ctx context.Context) (int, error)
	GetUsersRegisteredLastNDaysCount(ctx context.Context, since int) (int, error)

	GetRecentRegisteredUsers(ctx context.Context, limit int) ([]*models.User, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	UpdateUsername(ctx context.Context, id uuid.UUID, username string) error
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateAvatarUrl(ctx context.Context, id uuid.UUID, url string) error
	UpdateActive(ctx context.Context, id uuid.UUID, status bool) error
	UpdateVerify(ctx context.Context, id uuid.UUID, status bool) error
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error

	Search(ctx context.Context, search string) ([]*models.User, error)
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
