package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type ProblemHintRepository interface {
	Create(ctx context.Context, hint *models.ProblemHint) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetHints(ctx context.Context, problemId uuid.UUID) ([]*models.ProblemHint, error)
	Get(ctx context.Context, id uuid.UUID) (*models.ProblemHint, error)

	Update(ctx context.Context, id uuid.UUID, hint string) error
}

var (
	ErrProblemHintNotFound = errors.New("repository: hint not found")
)
