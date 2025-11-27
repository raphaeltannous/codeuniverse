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

type ProblemService interface {
	Create(ctx context.Context, problem *models.Problem) (*models.Problem, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetById(ctx context.Context, uuidString string) (*models.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*models.Problem, error)

	GetAllProblems(ctx context.Context, offset, limit int) ([]*models.Problem, error)

	UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error)
}

type problemService struct {
	problemRepository repository.ProblemRepository

	logger *slog.Logger
}

func NewProblemService(
	problemRepository repository.ProblemRepository,
) ProblemService {
	return &problemService{
		problemRepository: problemRepository,

		logger: slog.Default().With("package", "problemsService"),
	}
}

func (s *problemService) Create(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	uuid, err := s.problemRepository.Create(ctx, problem)
	if err != nil {
		s.logger.Error("failed to create problem", "problem", problem, "err", err)
		return nil, repository.ErrInternalServerError
	}

	problem.ID = uuid

	return problem, nil
}

func (s *problemService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *problemService) GetById(ctx context.Context, uuidString string) (*models.Problem, error) {
	return nil, nil
}

func (s *problemService) GetBySlug(ctx context.Context, slug string) (*models.Problem, error) {
	return nil, nil
}

func (s *problemService) GetAllProblems(ctx context.Context, offset, limit int) ([]*models.Problem, error) {
	return nil, nil
}

func (s *problemService) UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	return nil, nil
}
