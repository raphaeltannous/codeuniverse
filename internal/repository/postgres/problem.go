package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresProblemRepository struct {
	db *sql.DB
}

func NewProblemRepository(db *sql.DB) repository.ProblemRepository {
	return &postgresProblemRepository{db: db}
}

func (ppr *postgresProblemRepository) GetProblems(ctx context.Context, offset, limit int) ([]*models.Problem, error) {
	query := `
		SELECT
			id,

			title,
			slug,
			description,
			difficulty,

			to_json(hints) AS hints,

			code_snippets,
			test_cases,

			is_paid,
			is_public,

			created_at,
			updated_at
		FROM problems
		OFFSET $1
		LIMIT $2;
	`

	rows, err := ppr.db.QueryContext(
		ctx,
		query,
		offset,
		limit,
	)
	if err != nil {
		return nil, repository.ErrInternalServerError
	}
	defer rows.Close()

	var problems []*models.Problem
	for rows.Next() {
		problem := new(models.Problem)

		err := ppr.scanProblemFunc(rows, problem)
		if err != nil {
			return nil, repository.ErrInternalServerError
		}

		problems = append(problems, problem)
	}

	return problems, nil
}

func (ppr *postgresProblemRepository) Create(ctx context.Context, problem *models.Problem) (uuid.UUID, error) {
	query := `
		INSERT INTO problems (
			title,
			slug,
			description,
			difficulty,
			hints,
			code_snippets,
			test_cases,
			is_paid,
			is_public
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;
	`
	err := ppr.db.QueryRowContext(
		ctx,
		query,
		problem.Title,
		problem.Slug,
		problem.Description,
		problem.Difficulty,
		problem.Hints,
		problem.CodeSnippets,
		problem.TestCases,
		problem.IsPaid,
		problem.IsPublic,
	).Scan(
		&problem.ID,
	)

	return problem.ID, err
}

func (ppr *postgresProblemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (ppr *postgresProblemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Problem, error) {
	return nil, nil
}

func (ppr *postgresProblemRepository) GetBySlug(ctx context.Context, slug string) (*models.Problem, error) {
	return ppr.getProblemByColumn(
		ctx,
		"slug",
		slug,
	)
}

func (ppr *postgresProblemRepository) GetByNumber(ctx context.Context, number int) (*models.Problem, error) {
	return nil, nil
}

func (ppr *postgresProblemRepository) UpdateTitle(ctx context.Context, id uuid.UUID, title string) error {
	return ppr.updateColumnValue(
		ctx,
		id,
		"title",
		title,
	)
}

func (ppr *postgresProblemRepository) UpdateSlug(ctx context.Context, id uuid.UUID, slug string) error {
	return ppr.updateColumnValue(
		ctx,
		id,
		"slug",
		slug,
	)
}

func (ppr *postgresProblemRepository) UpdateDescription(ctx context.Context, id uuid.UUID, description string) error {
	return ppr.updateColumnValue(
		ctx,
		id,
		"description",
		description,
	)
}

func (ppr *postgresProblemRepository) UpdateDifficulty(ctx context.Context, id uuid.UUID, difficulty string) error {
	return ppr.updateColumnValue(
		ctx,
		id,
		"difficulty",
		difficulty,
	)
}

func (ppr *postgresProblemRepository) AddTags(ctx context.Context, id uuid.UUID, tags []string) error {
	return nil
}

func (ppr *postgresProblemRepository) AddTag(ctx context.Context, id uuid.UUID, tag string) error {
	return nil
}

func (ppr *postgresProblemRepository) RemoveTag(ctx context.Context, id uuid.UUID, tag string) error {
	return nil
}

func (ppr *postgresProblemRepository) AddHints(ctx context.Context, id uuid.UUID, hints []string) error {
	return nil
}

func (ppr *postgresProblemRepository) AddHint(ctx context.Context, id uuid.UUID, hint string) error {
	return nil
}

func (ppr *postgresProblemRepository) RemoveHint(ctx context.Context, id uuid.UUID, hint string) error {
	return nil
}

func (ppr *postgresProblemRepository) UpdateCodeSnippets(ctx context.Context, id uuid.UUID, codeSnippets string) error {
	return nil
}

func (ppr *postgresProblemRepository) UpdateTestcases(ctx context.Context, id uuid.UUID, testCases string) error {
	return nil
}

func (ppr *postgresProblemRepository) UpdatePublic(ctx context.Context, id uuid.UUID, status bool) error {
	return ppr.updateColumnValue(
		ctx,
		id,
		"is_public",
		status,
	)
}

func (ppr *postgresProblemRepository) UpdatePaid(ctx context.Context, id uuid.UUID, status bool) error {
	return ppr.updateColumnValue(
		ctx,
		id,
		"is_paid",
		status,
	)
}

func (ppr *postgresProblemRepository) updateColumnValue(ctx context.Context, id uuid.UUID, column string, value any) error {
	return updateColumnValue(
		ctx,
		ppr.db,
		"problems",
		id,
		column,
		value,
	)
}

func (ppr *postgresProblemRepository) Search(ctx context.Context, title string) ([]*models.Problem, error) {
	query := `
		SELECT
			id,

			title,
			slug,
			description,
			difficulty,

			to_json(hints) AS hints,

			code_snippets,
			test_cases,

			is_paid,
			is_public,

			created_at,
			updated_at
		FROM problems
		WHERE
			title ILIKE '%' || $1 || '%'
			OR slug ILIKE '%' || $1 || '%';
	`

	rows, err := ppr.db.QueryContext(
		ctx,
		query,
		title,
	)
	if err != nil {
		return nil, repository.ErrInternalServerError
	}
	defer rows.Close()

	var problems []*models.Problem
	for rows.Next() {
		problem := new(models.Problem)

		err := ppr.scanProblemFunc(rows, problem)
		if err != nil {
			return nil, repository.ErrInternalServerError
		}

		problems = append(problems, problem)
	}

	return problems, nil
}

func (ppr *postgresProblemRepository) getProblemByColumn(ctx context.Context, column string, value any) (*models.Problem, error) {
	query := fmt.Sprintf(
		`
		SELECT
			id,

			title,
			slug,
			description,
			difficulty,

			to_json(hints) AS hints,

			code_snippets,
			test_cases,

			is_paid,
			is_public,

			created_at,
			updated_at
		FROM problems
		WHERE %s = $1;
		`,
		column,
	)

	row := ppr.db.QueryRowContext(ctx, query, value)

	problem := new(models.Problem)
	if err := ppr.scanProblemFunc(row, problem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProblemNotFound
		}

		return nil, err
	}

	return problem, nil
}

type problemScanner interface {
	Scan(dest ...any) error
}

func (ppr *postgresProblemRepository) scanProblemFunc(scanner problemScanner, problem *models.Problem) error {
	var hintsBytes []byte
	var codeSnippetsBytes []byte
	var testCasesBytes []byte

	err := scanner.Scan(
		&problem.ID,

		&problem.Title,
		&problem.Slug,
		&problem.Description,
		&problem.Difficulty,

		&hintsBytes,

		&codeSnippetsBytes,
		&testCasesBytes,

		&problem.IsPaid,
		&problem.IsPublic,

		&problem.CreatedAt,
		&problem.UpdatedAt,
	)
	if err != nil {
		return err
	}

	var hints []string

	if len(hintsBytes) == 0 || string(hintsBytes) == "null" {
		hints = []string{}
	} else {
		if err := json.Unmarshal(hintsBytes, &hints); err != nil {
			return err
		}
	}

	problem.Hints = hints

	var codeSnippets []models.CodeSnippet

	if len(codeSnippetsBytes) == 0 || string(hintsBytes) == "null" {
		codeSnippets = []models.CodeSnippet{}
	} else {
		if err := json.Unmarshal(codeSnippetsBytes, &codeSnippets); err != nil {
			return err
		}
	}

	problem.CodeSnippets = codeSnippets

	var testCases []string

	if len(testCasesBytes) == 0 || string(testCasesBytes) == "null" {
		testCases = []string{}
	} else {
		if err := json.Unmarshal(testCasesBytes, &testCases); err != nil {
			return err
		}
	}

	problem.TestCases = testCases

	return nil
}
