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

type postgresProblemRepository struct {
	db *sql.DB
}

func NewProblemRepository(db *sql.DB) repository.ProblemRepository {
	return &postgresProblemRepository{db: db}
}

func (p *postgresProblemRepository) GetProblems(
	ctx context.Context,
	params *repository.GetProblemsParams,
) ([]*models.Problem, int, error) {
	whereClauses := []string{"1 = 1"}
	arguments := make([]any, 0)
	argumentPosition := 1

	if params.Search != "" {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"(title ILIKE '%%' || $%d || '%%' OR slug ILIKE '%%' || $%d || '%%')",
				argumentPosition,
				argumentPosition,
			),
		)
		arguments = append(arguments, params.Search)
		argumentPosition++
	}

	if params.IsPublic != 0 {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"is_public = $%d",
				argumentPosition,
			),
		)
		arguments = append(arguments, params.IsPublic == repository.ProblemPublic)
		argumentPosition++
	}

	if params.IsPremium != 0 {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(
				"is_premium = $%d",
				argumentPosition,
			),
		)
		arguments = append(arguments, params.IsPremium == repository.ProblemPremium)
		argumentPosition++
	}

	var orderBy strings.Builder
	switch params.SortBy {
	case repository.ProblemSortByTitle:
		orderBy.WriteString("title")
	default:
		orderBy.WriteString("created_at")
	}
	orderBy.WriteRune(' ')

	switch params.SortOrder {
	case repository.ProblemSortOrderAsc:
		orderBy.WriteString("ASC")
	default:
		orderBy.WriteString("DESC")
	}

	whereClause := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(
		`
			SELECT COUNT(*)
			FROM problems
			WHERE %s;
		`,
		whereClause,
	)

	var total int
	err := p.db.QueryRowContext(
		ctx,
		countQuery,
		arguments...,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count problems: %w", err)
	}

	query := fmt.Sprintf(
		`
			SELECT
				id,

				title,
				slug,
				description,
				difficulty,

				is_premium,
				is_public,

				created_at,
				updated_at
			FROM problems
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

	rows, err := p.db.QueryContext(
		ctx,
		query,
		arguments...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query problems: %w", err)
	}
	defer rows.Close()

	var problems []*models.Problem
	for rows.Next() {
		problem := new(models.Problem)

		err := p.scanProblemFunc(rows, problem)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan into problem: %w", err)
		}

		problems = append(problems, problem)
	}

	return problems, total, nil
}

func (p *postgresProblemRepository) Create(
	ctx context.Context,
	problem *models.Problem,
) (*models.Problem, error) {
	query := `
		INSERT INTO problems (
			title,
			slug,
			description,
			difficulty,

			is_premium,
			is_public
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	err := p.db.QueryRowContext(
		ctx,
		query,
		problem.Title,
		problem.Slug,
		problem.Description,
		problem.Difficulty,
		problem.IsPremium,
		problem.IsPublic,
	).Scan(
		&problem.ID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, repository.ErrProblemAlreadyExists
		}

		return nil, fmt.Errorf("failed to create new problem: %w", err)
	}

	return problem, nil
}

func (p *postgresProblemRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}

func (p *postgresProblemRepository) GetBySlug(
	ctx context.Context,
	slug string,
) (*models.Problem, error) {
	return p.getProblemByColumn(
		ctx,
		"slug",
		slug,
	)
}

func (p *postgresProblemRepository) GetCountByDifficulty(
	ctx context.Context,
	difficulty models.ProblemDifficulty,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM problems
		WHERE difficulty = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		difficulty.String(),
	)

	var difficultyCount int
	err := row.Scan(&difficultyCount)

	if err != nil {
		return 0, fmt.Errorf("failed to query count for difficulty (%s): %w", difficulty, err)
	}

	return difficultyCount, nil
}

func (p *postgresProblemRepository) UpdateTitle(
	ctx context.Context,
	id uuid.UUID,
	title string,
) error {
	return p.updateColumnValue(
		ctx,
		id,
		"title",
		title,
	)
}

func (p *postgresProblemRepository) UpdateSlug(
	ctx context.Context,
	id uuid.UUID,
	slug string,
) error {
	return p.updateColumnValue(
		ctx,
		id,
		"slug",
		slug,
	)
}

func (p *postgresProblemRepository) UpdateDescription(
	ctx context.Context,
	id uuid.UUID,
	description string,
) error {
	return p.updateColumnValue(
		ctx,
		id,
		"description",
		description,
	)
}

func (p *postgresProblemRepository) UpdateDifficulty(
	ctx context.Context,
	id uuid.UUID,
	difficulty models.ProblemDifficulty,
) error {
	return p.updateColumnValue(
		ctx,
		id,
		"difficulty",
		difficulty.String(),
	)
}

func (p *postgresProblemRepository) UpdateIsPremium(
	ctx context.Context,
	id uuid.UUID,
	status bool,
) error {
	return p.updateColumnValue(
		ctx,
		id,
		"is_premium",
		status,
	)
}

func (p *postgresProblemRepository) UpdateIsPublic(
	ctx context.Context,
	id uuid.UUID,
	status bool,
) error {
	return p.updateColumnValue(
		ctx,
		id,
		"is_public",
		status,
	)
}

func (p *postgresProblemRepository) updateColumnValue(
	ctx context.Context,
	id uuid.UUID,
	column string,
	value any,
) error {
	return updateColumnValue(
		ctx,
		p.db,
		"problems",
		id,
		column,
		value,
	)
}

func (p *postgresProblemRepository) getProblemByColumn(
	ctx context.Context,
	column string,
	value any,
) (*models.Problem, error) {
	query := fmt.Sprintf(
		`
			SELECT
				id,

				title,
				slug,
				description,
				difficulty,

				is_premium,
				is_public,

				created_at,
				updated_at
			FROM problems
			WHERE %s = $1;
		`,
		column,
	)

	row := p.db.QueryRowContext(ctx, query, value)

	problem := new(models.Problem)
	if err := p.scanProblemFunc(row, problem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProblemNotFound
		}

		return nil, fmt.Errorf("failed to scan into problem: %w", err)
	}

	return problem, nil
}

func (p *postgresProblemRepository) scanProblemFunc(
	scanner postgresScanner,
	problem *models.Problem,
) error {
	var difficultyString string

	err := scanner.Scan(
		&problem.ID,

		&problem.Title,
		&problem.Slug,
		&problem.Description,
		&difficultyString,

		&problem.IsPremium,
		&problem.IsPublic,

		&problem.CreatedAt,
		&problem.UpdatedAt,
	)
	if err != nil {
		return err
	}

	difficulty, err := models.NewProblemDifficulty(difficultyString)
	if err != nil {
		return err
	}
	problem.Difficulty = difficulty

	return nil
}
