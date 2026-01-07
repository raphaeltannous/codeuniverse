package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
)

type ProblemCodeRepository interface {
	DeleteProblem(ctx context.Context, problem *models.Problem) error

	SaveProblemCode(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode) error
	GetProblemCodes(ctx context.Context, problem *models.Problem) ([]*models.ProblemCode, error)

	SaveTestcases(ctx context.Context, problem *models.Problem, problemTestcases []*models.ProblemTestcase) error
	GetTestcases(ctx context.Context, problem *models.Problem) ([]*models.ProblemTestcase, error)

	SaveProblemCodeConfig(ctx context.Context, problem *models.Problem, config *models.ProblemCodeConfig) error
	GetProblemCodeConfig(ctx context.Context, problem *models.Problem) (*models.ProblemCodeConfig, error)
}

var (
	ErrCodeSnippetNotFound = errors.New("repository: snippet not found")
)
