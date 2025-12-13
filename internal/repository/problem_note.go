package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type ProblemNoteRepository interface {
	Create(ctx context.Context, note *models.ProblemNote) (*models.ProblemNote, error)
	Delete(ctx context.Context, id uuid.UUID) error

	UpdateMarkdown(ctx context.Context, note *models.ProblemNote, markdown string) error

	Get(ctx context.Context, userId, problemId uuid.UUID) (*models.ProblemNote, error)
}

var (
	ErrProblemNoteNotFound      = errors.New("repository: problem note not found")
	ErrProblemNoteAlreadyExists = errors.New("repository: problem note already exists")
)
