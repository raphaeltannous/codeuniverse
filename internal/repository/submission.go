package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) (*models.Submission, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetProblemSubmissions(ctx context.Context, userId, problemId uuid.UUID) ([]*models.Submission, error)

	GetById(ctx context.Context, id uuid.UUID) (*models.Submission, error)
}

var (
	ErrSubmissionNotFound = errors.New("repository: submission not found")
)
