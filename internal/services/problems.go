package services

import (
	"context"
	"errors"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidSlug = errors.New("invalid slug")
)

type ProblemsService interface {
	Create(ctx context.Context, problem *models.Problem) (*models.Problem, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetById(ctx context.Context, uuidString string) (*models.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*models.Problem, error)

	GetAllProblems(ctx context.Context, offset, limit int) ([]*models.Problem, error)

	UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error)
}

type problemsService struct {
	problemRepository repository.ProblemRepository

	logger *slog.Logger
}

func NewProblemService(
	problemRepository repository.ProblemRepository,
) ProblemsService {
	return &problemsService{
		problemRepository: problemRepository,

		logger: slog.Default().With("package", "problemsService"),
	}
}

func (s *problemsService) Create(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	return nil, nil
}

func (s *problemsService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *problemsService) GetById(ctx context.Context, uuidString string) (*models.Problem, error) {
	return nil, nil
}

func (s *problemsService) GetBySlug(ctx context.Context, slug string) (*models.Problem, error) {
	return nil, nil
}

func (s *problemsService) GetAllProblems(ctx context.Context, offset, limit int) ([]*models.Problem, error) {
	return nil, nil
}

func (s *problemsService) UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	return nil, nil
}
