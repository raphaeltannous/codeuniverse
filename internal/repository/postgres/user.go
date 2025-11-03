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

type postgresUserRepository struct {
	db *sql.DB
}

var _ repository.UserRepository = (*postgresUserRepository)(nil)

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (pur *postgresUserRepository) GetUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_verified, is_active, role, created_at, updated_at
		FROM users
		OFFSET $1
		LIMIT $2
	`

	rows, err := pur.db.QueryContext(
		ctx,
		query,
		offset,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)

		err := scanUserFunc(rows, user)
		if err != nil {
			return nil, fmt.Errorf("failed to scan into user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
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

func (pur *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (pur *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, email, is_verified, is_active, role, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	row := pur.db.QueryRowContext(ctx, query, id)

	user := new(models.User)
	if err := scanUser(row, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan user data into *model.User: %w", err)
	}

	return user, nil
}

func (pur *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}

func (pur *postgresUserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	return nil
}

func (pur *postgresUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	return nil
}

func (pur *postgresUserRepository) UpdateActive(ctx context.Context, id uuid.UUID, status bool) error {
	return nil
}

func (pur *postgresUserRepository) UpdateVerify(ctx context.Context, id uuid.UUID, status bool) error {
	return nil
}

func (pur *postgresUserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	return nil
}

func scanUser(row *sql.Row, user *models.User) error {
	return row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.IsVerified,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUserFunc(scanner userScanner, user *models.User) error {
	return scanner.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
