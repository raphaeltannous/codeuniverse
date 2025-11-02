package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresUserRepository struct {
	db *sql.DB
}

var _ repository.UserRepository = (*postgresUserRepository)(nil)

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (pur *postgresUserRepository) Create(ctx context.Context, user *models.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (username, password_hash, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	row := pur.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Role,
	)

	err := row.Scan(&user.ID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error inserting user: %w", err)
	}

	return user.ID, nil
}

func (pur *postgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) Update(ctx context.Context, u *models.User) (*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (pur *postgresUserRepository) VerifyEmail(ctx context.Context, id string) error {
	return nil
}

func (pur *postgresUserRepository) SetActive(ctx context.Context, id string, active bool) error {
	return nil
}
