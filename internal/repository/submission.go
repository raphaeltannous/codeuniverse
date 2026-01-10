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
	UpdateFailedTestcases(ctx context.Context, id uuid.UUID, failedTestcases []*models.FailedTestcase) error
	UpdateStdout(ctx context.Context, id uuid.UUID, stdout string) error
	UpdateStderr(ctx context.Context, id uuid.UUID, stderr string) error

	GetProblemSubmissions(ctx context.Context, userId, problemId uuid.UUID) ([]*models.Submission, error)
	GetSubmissionsStats(ctx context.Context, userId uuid.UUID) (*models.SubmissionStats, error)

	GetSubmissionsCount(ctx context.Context) (int, error)
	GetSubmissionsLastNDaysCount(ctx context.Context, since int) (int, error)
	GetPendingSubmissionsCount(ctx context.Context) (int, error)
	GetAcceptedSubmissionsCount(ctx context.Context) (int, error)
	GetRecentSubmissions(ctx context.Context, limit int) ([]*models.SubmissionActivity, error)
	GetDailySubmissions(ctx context.Context, since int) ([]*models.DailySubmissions, error)
	GetDailySubmissionsHours(ctx context.Context, since int) ([]*models.DailySubmissions, error)

	GetById(ctx context.Context, id uuid.UUID) (*models.Submission, error)
}

var (
	ErrSubmissionNotFound = errors.New("repository: submission not found")
)
