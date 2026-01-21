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

	GetUsers(ctx context.Context, params *GetUsersParams) ([]*models.User, int, error)

	GetAdminCount(ctx context.Context) (int, error)
	GetUsersCount(ctx context.Context) (int, error)
	GetUsersRegisteredLastNDaysCount(ctx context.Context, since int) (int, error)

	GetRecentRegisteredUsers(ctx context.Context, limit int) ([]*models.User, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByStripeCustomerId(ctx context.Context, customerId string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	UpdateUsername(ctx context.Context, id uuid.UUID, username string) error
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateAvatarUrl(ctx context.Context, id uuid.UUID, url string) error
	UpdateActive(ctx context.Context, id uuid.UUID, status bool) error
	UpdateVerify(ctx context.Context, id uuid.UUID, status bool) error
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error
	UpdateStripeCustomerId(ctx context.Context, id uuid.UUID, customerId string) error
	UpdatePremiumStatus(ctx context.Context, id uuid.UUID, status string) error
}

type UserParam int

const (
	UserInactive UserParam = iota + 1
	UserActive

	UserUnverified
	UserVerified

	UserRoleAdmin
	UserRoleUser

	UserSortByUsername
	UserSortByEmail
	UserSortByCreatedAt

	UserSortOrderAsc
	UserSortOrderDesc
)

type GetUsersParams struct {
	Offset int
	Limit  int
	Search string

	Role UserParam

	IsActive   UserParam
	IsVerified UserParam

	SortBy    UserParam
	SortOrder UserParam
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
