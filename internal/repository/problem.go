package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type ProblemRepository interface {
	GetAllProblems(ctx context.Context, limit, offset int) ([]*models.Problem, error)

	Create(ctx context.Context, problem *models.Problem) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*models.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*models.Problem, error)

	UpdateTitle(ctx context.Context, title string) error
	UpdateSlug(ctx context.Context, slug string) error

	UpdateDescription(ctx context.Context, descripton string) error

	UpdateDifficulty(ctx context.Context, difficulty string) error

	AddTags(ctx context.Context, tags []string) error
	AddTag(ctx context.Context, tag string) error
	RemoveTag(ctx context.Context, tag string) error

	AddHints(ctx context.Context, hints []string) error
	AddHint(ctx context.Context, hint string) error
	RemoveHint(ctx context.Context, hint string) error

	UpdateCodeSnippets(ctx context.Context, codeSnippets string) error
	UpdateTestcases(ctx context.Context, testCases string) error

	UpdatePublic(ctx context.Context, id uuid.UUID, status bool) error
	UpdatePaind(ctx context.Context, id uuid.UUID, status bool) error

	IncrementLikes(ctx context.Context, id uuid.UUID, incrementValue int) error
	DecrementLikes(ctx context.Context, id uuid.UUID, decrementValue int) error
	IncrementDislikes(ctx context.Context, id uuid.UUID, incrementValue int) error
	DecrementDislikes(ctx context.Context, id uuid.UUID, decrementValue int) error
}

var (
	ErrProblemNotFound      = errors.New("problem not found")
	ErrProblemAlreadyExists = errors.New("problem already exists")
)
