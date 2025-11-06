package postgres

import (
	"context"
	"database/sql"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresProblemRepository struct {
	db *sql.DB
}

var _ repository.ProblemRepository = (*postgresProblemRepository)(nil)

func (ppr *postgresProblemRepository) GetProblems(ctx context.Context, limit, offset int) ([]*models.Problem, error) {
	return nil, nil
}

func (ppr *postgresProblemRepository) Create(ctx context.Context, problem *models.Problem) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (ppr *postgresProblemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (ppr *postgresProblemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Problem, error) {
	return nil, nil
}

func (ppr *postgresProblemRepository) GetBySlug(ctx context.Context, slug string) (*models.Problem, error) {
	return nil, nil
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
