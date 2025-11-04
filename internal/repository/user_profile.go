package repository

import (
	"context"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type UserProfileRepository interface {
	GetInfo(ctx context.Context, userId uuid.UUID) (*models.UserProfile, error)

	UpdateName(ctx context.Context, userId uuid.UUID, name string) error
	UpdateBio(ctx context.Context, userId uuid.UUID, bio string) error
	UpdateAvatarURL(ctx context.Context, userId uuid.UUID, url string) error
	UpdateCountry(ctx context.Context, userId uuid.UUID, country string) error
	UpdatePreferredLanguage(ctx context.Context, userId uuid.UUID, language string) error

	UpdateWebsiteURL(ctx context.Context, userId uuid.UUID, url string) error
	UpdateGithubURL(ctx context.Context, userId uuid.UUID, url string) error
	UpdateLinkedinURL(ctx context.Context, userId uuid.UUID, url string) error
	UpdateXURL(ctx context.Context, userId uuid.UUID, url string) error

	UpdateLastActive(ctx context.Context, userId uuid.UUID) error
}
