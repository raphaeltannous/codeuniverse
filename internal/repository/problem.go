package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type ProblemRepository interface {
	GetProblems(ctx context.Context, limit, offset int) ([]*models.Problem, error)

	Create(ctx context.Context, problem *models.Problem) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*models.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*models.Problem, error)
	GetByNumber(ctx context.Context, number int) (*models.Problem, error)

	UpdateTitle(ctx context.Context, id uuid.UUID, title string) error
	UpdateSlug(ctx context.Context, id uuid.UUID, slug string) error

	UpdateDescription(ctx context.Context, id uuid.UUID, description string) error

	UpdateDifficulty(ctx context.Context, id uuid.UUID, difficulty string) error

	AddTags(ctx context.Context, id uuid.UUID, tags []string) error
	AddTag(ctx context.Context, id uuid.UUID, tag string) error
	RemoveTag(ctx context.Context, id uuid.UUID, tag string) error

	AddHints(ctx context.Context, id uuid.UUID, hints []string) error
	AddHint(ctx context.Context, id uuid.UUID, hint string) error
	RemoveHint(ctx context.Context, id uuid.UUID, hint string) error

	UpdateCodeSnippets(ctx context.Context, id uuid.UUID, codeSnippets string) error
	UpdateTestcases(ctx context.Context, id uuid.UUID, testCases string) error

	UpdatePublic(ctx context.Context, id uuid.UUID, status bool) error
	UpdatePaid(ctx context.Context, id uuid.UUID, status bool) error

	// TODO:
	// - Search capabilities.
	// - Get all problems?
}

var (
	ErrProblemNotFound      = errors.New("problem not found")
	ErrProblemAlreadyExists = errors.New("problem already exists")
)
