package postgres

import (
	"context"
	"database/sql"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresUserProfileRepository struct {
	db *sql.DB
}

var _ repository.UserProfileRepository = (*postgresUserProfileRepository)(nil)

func (pupr *postgresUserProfileRepository) Create(ctx context.Context, userId uuid.UUID) error {
	// TODO
	return nil
}

func (pupr *postgresUserProfileRepository) GetInfo(ctx context.Context, userId uuid.UUID) (*models.UserProfile, error) {
	// TODO
	return nil, nil
}

func (pupr *postgresUserProfileRepository) UpdateName(ctx context.Context, userId uuid.UUID, name string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"name",
		name,
	)
}

func (pupr *postgresUserProfileRepository) UpdateBio(ctx context.Context, userId uuid.UUID, bio string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"bio",
		bio,
	)
}

func (pupr *postgresUserProfileRepository) UpdateAvatarURL(ctx context.Context, userId uuid.UUID, url string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"avatar_url",
		url,
	)
}

func (pupr *postgresUserProfileRepository) UpdateCountry(ctx context.Context, userId uuid.UUID, country string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"country",
		country,
	)
}

func (pupr *postgresUserProfileRepository) UpdatePreferredLanguage(ctx context.Context, userId uuid.UUID, language string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"preferred_language",
		language,
	)
}

func (pupr *postgresUserProfileRepository) UpdateWebsiteURL(ctx context.Context, userId uuid.UUID, url string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"website_url",
		url,
	)
}

func (pupr *postgresUserProfileRepository) UpdateGithubURL(ctx context.Context, userId uuid.UUID, url string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"github_url",
		url,
	)
}

func (pupr *postgresUserProfileRepository) UpdateLinkedinURL(ctx context.Context, userId uuid.UUID, url string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"linkedin_url",
		url,
	)
}

func (pupr *postgresUserProfileRepository) UpdateXURL(ctx context.Context, userId uuid.UUID, url string) error {
	return pupr.updateColumnValue(
		ctx,
		userId,
		"x_url",
		url,
	)
}

func (pupr *postgresUserProfileRepository) UpdateLastActive(ctx context.Context, userId uuid.UUID) error {
	// TODO
	return nil
}

func (pupr *postgresUserProfileRepository) updateColumnValue(ctx context.Context, id uuid.UUID, column string, value any) error {
	return updateColumnValue(
		ctx,
		pupr.db,
		"user_profiles",
		id,
		column,
		value,
	)
}
