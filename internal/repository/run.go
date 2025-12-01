package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type RunRepository interface {
	Create(ctx context.Context, run *models.Run) (*models.Run, error)

	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateExecutionTime(ctx context.Context, id uuid.UUID, executionTime float64) error
	UpdateMemoryUsage(ctx context.Context, id uuid.UUID, memoryUsage float64) error
	UpdateAcceptanceStatus(ctx context.Context, id uuid.UUID, status bool) error

	GetById(ctx context.Context, id uuid.UUID) (*models.Run, error)
}

var (
	ErrRunNotFound = errors.New("repository: run not found")
)
