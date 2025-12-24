package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresUserProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) repository.UserProfileRepository {
	return &postgresUserProfileRepository{db: db}
}

func (pupr *postgresUserProfileRepository) Create(ctx context.Context, userId uuid.UUID) error {
	query := `
		INSERT INTO user_profiles
			(user_id)
		VALUES ($1);
	`

	// TODO: should I do something about sql.Result?
	_, err := getExecutor(ctx, pupr.db).ExecContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	return nil
}

func (pupr *postgresUserProfileRepository) GetInfo(ctx context.Context, user *models.User) (*models.UserProfile, error) {
	query := `
		SELECT
			user_id,

			name,
			bio,
			avatar_url,
			country,

			preferred_language,

			website_url,
			github_url,
			linkedin_url,
			x_url,

			last_active,
			created_at,
			updated_at
		FROM user_profiles
		WHERE user_id = $1;
	`

	row := pupr.db.QueryRowContext(
		ctx,
		query,
		user.ID,
	)

	userProfile := new(models.UserProfile)
	if err := pupr.scanUserProfileFunc(row, userProfile); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserProfileNotFound
		}

		return nil, fmt.Errorf("failed to scan userProfile into *model.UserProfile: %w", err)
	}

	return userProfile, nil
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

func (pupr *postgresUserProfileRepository) scanUserProfileFunc(scanner postgresScanner, userProfile *models.UserProfile) error {
	return scanner.Scan(
		&userProfile.UserID,
		&userProfile.Name,
		&userProfile.Bio,
		&userProfile.AvatarURL,
		&userProfile.Country,
		&userProfile.PreferredLanguage,
		&userProfile.WebsiteURL,
		&userProfile.GithubURL,
		&userProfile.LinkedinURL,
		&userProfile.XURL,
		&userProfile.LastActive,
		&userProfile.CreatedAt,
		&userProfile.UpdatedAt,
	)
}
