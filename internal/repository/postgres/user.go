package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (pur *postgresUserRepository) GetUsers(ctx context.Context, params *repository.GetUsersParams) ([]*models.User, int, error) {
	whereClauses := []string{"1 = 1"}
	arguments := make([]any, 0)
	argumentPosition := 1

	if params.Search != "" {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"(username ILIKE '%%' || $%d || '%%' OR email ILIKE '%%' || $%d || '%%')",
				argumentPosition,
				argumentPosition,
			),
		)
		arguments = append(arguments, params.Search)
		argumentPosition++
	}

	if params.Role != "" {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"role = $%d",
				argumentPosition,
			),
		)
		arguments = append(arguments, params.Role)
		argumentPosition++
	}

	if params.IsActive != 0 {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"is_active = $%d",
				argumentPosition,
			),
		)
		arguments = append(arguments, params.IsActive == repository.UserActive)
		argumentPosition++
	}

	if params.IsVerified != 0 {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"is_verified = $%d",
				argumentPosition,
			),
		)
		arguments = append(arguments, params.IsVerified == repository.UserVerified)
		argumentPosition++
	}

	var orderBy strings.Builder
	switch params.SortBy {
	case repository.UserSortByUsername:
		orderBy.WriteString("username")
	case repository.UserSortByEmail:
		orderBy.WriteString("email")
	default:
		orderBy.WriteString("created_at")
	}
	orderBy.WriteString(" ")

	switch params.SortOrder {
	case repository.UserSortOrderAsc:
		orderBy.WriteString("ASC")
	default:
		orderBy.WriteString("DESC")
	}

	whereClause := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(
		`
			SELECT COUNT(*)
			FROM users
			WHERE %s;
		`,
		whereClause,
	)

	var total int
	err := pur.db.QueryRowContext(
		ctx,
		countQuery,
		arguments...,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	query := fmt.Sprintf(
		`
			SELECT
				id,

				username,
				email,
				password_hash,
				avatar_url,

				is_verified,
				is_active,
				role,

				created_at,
				updated_at
			FROM users
			WHERE %s
			ORDER BY %s
			OFFSET $%d
			LIMIT $%d;
		`,
		whereClause,
		orderBy.String(),
		argumentPosition,
		argumentPosition+1,
	)
	argumentPosition++

	arguments = append(arguments, params.Offset, params.Limit)

	rows, err := pur.db.QueryContext(
		ctx,
		query,
		arguments...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)

		err := scanUserFunc(rows, user)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, user)
	}

	return users, total, nil
}

func (pur *postgresUserRepository) GetRecentRegisteredUsers(ctx context.Context, limit int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, avatar_url, is_verified, is_active, role, created_at, updated_at
		FROM users
		WHERE created_at >= NOW() - INTERVAL '24 hours'
		ORDER BY created_at DESC
		LIMIT $1;
	`

	rows, err := pur.db.QueryContext(
		ctx,
		query,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent registered users: %w", err)
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

func (pur *postgresUserRepository) GetAdminCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE role = 'admin';
	`

	row := pur.db.QueryRowContext(
		ctx,
		query,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get admin count: %w", err)
	}

	return count, nil
}

func (pur *postgresUserRepository) GetUsersCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users;
	`

	row := pur.db.QueryRowContext(
		ctx,
		query,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get users count: %w", err)
	}

	return count, nil
}

func (pur *postgresUserRepository) GetUsersRegisteredLastNDaysCount(ctx context.Context, since int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE created_at >= NOW() - $1::INTERVAL;
	`

	row := pur.db.QueryRowContext(
		ctx,
		query,
		fmt.Sprintf("%d days", since),
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get users count since %d days: %w", since, err)
	}

	return count, nil
}

func (pur *postgresUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	row := getExecutor(ctx, pur.db).QueryRowContext(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Role,
	)

	err := row.Scan(&user.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, repository.ErrUserAlreadyExists
		}

		return user, fmt.Errorf("error inserting user: %w", err)
	}

	return user, nil
}

func (pur *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		where id = $1;
	`

	result, err := pur.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete requested user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to query rows affected from database: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("expect single row affected, got %d rows affected", rows)
	}

	return nil
}

func (pur *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return pur.getUserByColumn(
		ctx,
		"id",
		id,
	)
}

func (pur *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return pur.getUserByColumn(
		ctx,
		"username",
		username,
	)
}

func (pur *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return pur.getUserByColumn(
		ctx,
		"email",
		email,
	)
}

func (pur *postgresUserRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"username",
		username,
	)
}

func (pur *postgresUserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"email",
		email,
	)
}

func (pur *postgresUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"password_hash",
		password,
	)
}

func (pur *postgresUserRepository) UpdateAvatarUrl(ctx context.Context, id uuid.UUID, url string) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"avatar_url",
		url,
	)
}

func (pur *postgresUserRepository) UpdateActive(ctx context.Context, id uuid.UUID, status bool) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"is_active",
		status,
	)
}

func (pur *postgresUserRepository) UpdateVerify(ctx context.Context, id uuid.UUID, status bool) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"is_verified",
		status,
	)
}

func (pur *postgresUserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	return pur.updateColumnValue(
		ctx,
		id,
		"role",
		role,
	)
}

func (pur *postgresUserRepository) getUserByColumn(ctx context.Context, column string, value any) (*models.User, error) {
	query := fmt.Sprintf(
		`
			SELECT id, username, email, password_hash, avatar_url, is_verified, is_active, role, created_at, updated_at
			FROM users
			WHERE %s = $1;
		`,
		column,
	)

	row := pur.db.QueryRowContext(ctx, query, value)

	user := new(models.User)
	// TODO: what are the erros that are returnes by row.Scan()?
	if err := scanUserFunc(row, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan user data into *model.User: %w", err)
	}

	return user, nil
}

func (pur *postgresUserRepository) Search(ctx context.Context, search string) ([]*models.User, error) {
	query := `
		SELECT
			id,

			username,
			email,
			password_hash,

			is_verified,
			is_active,

			role,

			created_at,
			updated_at
		FROM users
		WHERE
			username ILIKE '%' || $1 || '%'
			OR email ILIKE '%' || $1 || '%';
	`

	rows, err := pur.db.QueryContext(
		ctx,
		query,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query users by search: %w", err)
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

func (pur *postgresUserRepository) updateColumnValue(ctx context.Context, id uuid.UUID, column string, value any) error {
	return updateColumnValue(
		ctx,
		pur.db,
		"users",
		id,
		column,
		value,
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
		&user.AvatarURL,
		&user.IsVerified,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
