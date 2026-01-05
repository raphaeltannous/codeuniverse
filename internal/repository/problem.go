package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type ProblemRepository interface {
	GetProblems(ctx context.Context, params *GetProblemsParams) ([]*models.Problem, int, error)

	Create(ctx context.Context, problem *models.Problem) (*models.Problem, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetBySlug(ctx context.Context, slug string) (*models.Problem, error)

	GetCountByDifficulty(ctx context.Context, difficulty models.ProblemDifficulty) (int, error)

	UpdateTitle(ctx context.Context, id uuid.UUID, title string) error
	UpdateSlug(ctx context.Context, id uuid.UUID, slug string) error
	UpdateDescription(ctx context.Context, id uuid.UUID, description string) error
	UpdateDifficulty(ctx context.Context, id uuid.UUID, difficulty models.ProblemDifficulty) error
	UpdateIsPremium(ctx context.Context, id uuid.UUID, status bool) error
	UpdateIsPublic(ctx context.Context, id uuid.UUID, status bool) error
}

type ProblemParam int

const (
	ProblemPublic ProblemParam = iota + 1
	ProblemPrivate

	ProblemFree
	ProblemPremium

	ProblemSortByTitle
	ProblemSortByCreatedAt

	ProblemSortOrderAsc
	ProblemSortOrderDesc
)

type GetProblemsParams struct {
	Offset int
	Limit  int
	Search string

	IsPublic  ProblemParam
	IsPremium ProblemParam

	Difficulty models.ProblemDifficulty

	SortBy    ProblemParam
	SortOrder ProblemParam
}

var (
	ErrProblemAlreadyExists = errors.New("repository: problem already exists")
	ErrProblemNotFound      = errors.New("repository: problem not found")
)
