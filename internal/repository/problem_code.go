package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
)

type ProblemCodeRepository interface {
	SaveCodeSnippet(ctx context.Context, problem *models.Problem, snippet *models.CodeSnippet) error
	// GetCodeSnippets(ctx context.Context, problem *models.Problem) ([]*models.CodeSnippet, error)
	// DeleteCodeSnippet(ctx context.Context, problem *models.Problem, language *models.ProblemLanguage) error

	// SaveDriverCode(ctx context.Context, problem *models.Problem, code string) error
	// GetDriverCode(ctx context.Context, problem *models.Problem) (string, error)
	// DeleteDriverCode(ctx context.Context, problem *models.Problem) error

	// SaveSolutionCode(ctx context.Context, problem *models.Problem, code string) error
	// GetSolutionCode(ctx context.Context, problem *models.Problem) (string, error)
	// DeleteSolutionCode(ctx context.Context, problem *models.Problem) error

	// SaveTests(ctx context.Context, problem *models.Problem, problemTestcases *models.ProblemTestcases) error
	// GetTests(ctx context.Context, problem *models.Problem) (*models.ProblemTestcases, error)
	// DeleteTestcases(ctx context.Context, problem *models.Problem) error
}

var (
	ErrCodeSnippetNotFound = errors.New("repository: snippet not found")
)
