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

	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateExecutionTime(ctx context.Context, id uuid.UUID, executionTime float64) error
	UpdateMemoryUsage(ctx context.Context, id uuid.UUID, memoryUsage float64) error
	UpdateAcceptanceStatus(ctx context.Context, id uuid.UUID, status bool) error

	GetProblemSubmissions(ctx context.Context, userId, problemId uuid.UUID) ([]*models.Submission, error)
	GetSubmissionsStats(ctx context.Context, userId uuid.UUID) (*models.SubmissionStats, error)

	GetById(ctx context.Context, id uuid.UUID) (*models.Submission, error)
}

var (
	ErrSubmissionNotFound = errors.New("repository: submission not found")
)
