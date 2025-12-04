package services

import (
	"context"
	"errors"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/judger"
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

	Submit(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error
	Run(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error
}

type problemService struct {
	problemRepository    repository.ProblemRepository
	runRepository        repository.RunRepository
	submissionRepository repository.SubmissionRepository

	judge  judger.Judge
	logger *slog.Logger
}

func NewProblemService(
	problemRepository repository.ProblemRepository,
	runRepository repository.RunRepository,
	submissionRepository repository.SubmissionRepository,

	judge judger.Judge,
) ProblemService {
	return &problemService{
		problemRepository:    problemRepository,
		runRepository:        runRepository,
		submissionRepository: submissionRepository,

		judge:  judge,
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
	problem, err := s.problemRepository.GetBySlug(
		ctx,
		slug,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProblemNotFound):
			return nil, err
		}

		s.logger.Error("failed to get err", "func", "GetBySlug", "err", err)
		return nil, repository.ErrInternalServerError
	}

	// Filtering snippets based on the available snippets in judger.problemsDataDir
	allowedSnippets := make([]models.CodeSnippet, 0, len(problem.CodeSnippets))
	for _, snippet := range problem.CodeSnippets {
		if language, ok := judger.SupportedLanguages[snippet.LanguageSlug]; ok {
			if language.DoesItHaveTests(problem.Slug) {
				allowedSnippets = append(allowedSnippets, snippet)
			}
		}
	}
	problem.CodeSnippets = allowedSnippets

	return problem, nil
}

func (s *problemService) GetAllProblems(ctx context.Context, offset, limit int) ([]*models.Problem, error) {
	problems, err := s.problemRepository.GetProblems(
		ctx,
		offset,
		limit,
	)
	if err != nil {
		return nil, repository.ErrInternalServerError
	}
	return problems, nil
}

func (s *problemService) UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	return nil, nil
}

func (s *problemService) Submit(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error {
	submission := models.NewSubmission(
		user.ID,
		problem.ID,

		languageSlug,
		code,
		"PENDING",
	)

	submission, err := s.submissionRepository.Create(
		ctx,
		submission,
	)
	if err != nil {
		handlerChannel <- repository.ErrInternalServerError.Error()
		close(handlerChannel)
		return err
	}
	handlerChannel <- submission.ID.String()
	close(handlerChannel)

	err = s.judge.Submit(
		ctx,
		submission,
		problem.Slug,
	)
	if err != nil {
		s.logger.Error("failed to submit judge", "err", err)
		return err
	}

	if err := s.submissionRepository.UpdateAcceptanceStatus(ctx, submission.ID, submission.IsAccepted); err != nil {
		return err
	}

	if err := s.submissionRepository.UpdateStatus(ctx, submission.ID, submission.Status); err != nil {
		return err
	}

	return nil
}

func (s *problemService) Run(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error {
	run := models.NewRun(
		user.ID,
		problem.ID,
		languageSlug,
		code,
		"PENDING",
	)

	run, err := s.runRepository.Create(
		ctx,
		run,
	)
	if err != nil {
		handlerChannel <- repository.ErrInternalServerError.Error()
		return err
	}

	handlerChannel <- run.ID.String()
	close(handlerChannel)

	err = s.judge.Run(
		ctx,
		run,
		problem.Slug,
	)
	if err != nil {
		s.logger.Error("failed to run judge", "err", err)
		return err
	}

	if err := s.runRepository.UpdateAcceptanceStatus(ctx, run.ID, run.IsAccepted); err != nil {
		return err
	}

	if err := s.runRepository.UpdateStatus(ctx, run.ID, run.Status); err != nil {
		return err
	}

	return nil
}
