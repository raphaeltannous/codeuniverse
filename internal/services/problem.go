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

	GetEasyCount(ctx context.Context) (int, error)
	GetMediumCount(ctx context.Context) (int, error)
	GetHardCount(ctx context.Context) (int, error)

	UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error)

	Submit(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error
	Run(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error

	GetRun(ctx context.Context, user *models.User, runId uuid.UUID) (*models.Run, error)

	GetSubmission(ctx context.Context, user *models.User, submissionId uuid.UUID) (*models.Submission, error)
	GetSubmissions(ctx context.Context, user *models.User, problem *models.Problem) ([]*models.Submission, error)

	GetSubmissionsCount(ctx context.Context) (int, error)
	GetSubmissionsLastNDaysCount(ctx context.Context, since int) (int, error)
	GetPendingSubmissionsCount(ctx context.Context) (int, error)
	GetAcceptedSubmissionsCount(ctx context.Context) (int, error)
	GetRecentSubmissions(ctx context.Context, count int) ([]*models.SubmissionActivity, error)
	GetDailySubmissions(ctx context.Context, since int) ([]*models.DailySubmissions, error)

	CreateNote(ctx context.Context, note *models.ProblemNote) error
	DeleteNote(ctx context.Context, note *models.ProblemNote) error
	GetNote(ctx context.Context, user *models.User, problem *models.Problem) (*models.ProblemNote, error)
	UpdateNote(ctx context.Context, note *models.ProblemNote, markdown string) error

	Search(ctx context.Context, search string) ([]*models.Problem, error)
}

type problemService struct {
	problemRepository     repository.ProblemRepository
	problemNoteRepository repository.ProblemNoteRepository
	runRepository         repository.RunRepository
	submissionRepository  repository.SubmissionRepository

	judge  judger.Judge
	logger *slog.Logger
}

func NewProblemService(
	problemRepository repository.ProblemRepository,
	problemNoteRepository repository.ProblemNoteRepository,
	runRepository repository.RunRepository,
	submissionRepository repository.SubmissionRepository,

	judge judger.Judge,
) ProblemService {
	return &problemService{
		problemRepository:     problemRepository,
		problemNoteRepository: problemNoteRepository,
		runRepository:         runRepository,
		submissionRepository:  submissionRepository,

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
	s.logger.Debug("offset and limit", "offset", offset, "limit", limit)
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

func (s *problemService) GetEasyCount(ctx context.Context) (int, error) {
	count, err := s.problemRepository.GetEasyCount(ctx)
	if err != nil {
		s.logger.Error("failed to get easy count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) GetMediumCount(ctx context.Context) (int, error) {
	count, err := s.problemRepository.GetMediumCount(ctx)
	if err != nil {
		s.logger.Error("failed to get medium count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) GetHardCount(ctx context.Context) (int, error) {
	count, err := s.problemRepository.GetHardCount(ctx)
	if err != nil {
		s.logger.Error("failed to get hard count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
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

	if err := s.submissionRepository.UpdateExecutionTime(ctx, submission.ID, submission.ExecutionTime); err != nil {
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

func (s *problemService) GetRun(ctx context.Context, user *models.User, runId uuid.UUID) (*models.Run, error) {
	run, err := s.runRepository.GetById(
		ctx,
		runId,
	)
	if err != nil {
		if !errors.Is(err, repository.ErrRunNotFound) {
			s.logger.Error("repository: failed to get run", "err", err, "runId", runId)
		}
		return nil, err
	}

	if run.UserId != user.ID {
		s.logger.Error("User requesting run for another user", "user.ID", user.ID, "run.UserId", run.UserId)
		return nil, errors.New("Access denied.")
	}

	return run, nil
}

func (s *problemService) GetSubmission(ctx context.Context, user *models.User, submissionId uuid.UUID) (*models.Submission, error) {
	submission, err := s.submissionRepository.GetById(
		ctx,
		submissionId,
	)
	if err != nil {
		if !errors.Is(err, repository.ErrSubmissionNotFound) {
			s.logger.Error("repository: failed to get submission", "err", err, "submissionId", submissionId)
		}

		return nil, err
	}

	if submission.UserId != user.ID {
		s.logger.Error("User requesting submission for another user", "user.ID", user.ID, "submission.UserId", submission.UserId)
		return nil, errors.New("Access denied.") // TODO: universal Access denied?
	}

	return submission, nil
}

func (s *problemService) GetSubmissions(ctx context.Context, user *models.User, problem *models.Problem) ([]*models.Submission, error) {
	submissions, err := s.submissionRepository.GetProblemSubmissions(
		ctx,
		user.ID,
		problem.ID,
	)
	if err != nil {
		s.logger.Error("repository: get submissions", "err", err, "user", user, "problem", problem)
		return nil, err
	}

	return submissions, nil
}

func (s *problemService) GetSubmissionsCount(ctx context.Context) (int, error) {
	count, err := s.submissionRepository.GetSubmissionsCount(ctx)
	if err != nil {
		s.logger.Error("failed to get submissions count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) GetSubmissionsLastNDaysCount(ctx context.Context, since int) (int, error) {
	count, err := s.submissionRepository.GetSubmissionsLastNDaysCount(ctx, since)
	if err != nil {
		s.logger.Error("failed to get submissions count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) GetPendingSubmissionsCount(ctx context.Context) (int, error) {
	count, err := s.submissionRepository.GetPendingSubmissionsCount(ctx)
	if err != nil {
		s.logger.Error("failed to get submissions count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) GetAcceptedSubmissionsCount(ctx context.Context) (int, error) {
	count, err := s.submissionRepository.GetAcceptedSubmissionsCount(ctx)
	if err != nil {
		s.logger.Error("failed to get submissions count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) GetRecentSubmissions(ctx context.Context, count int) ([]*models.SubmissionActivity, error) {
	submissionsActivities, err := s.submissionRepository.GetRecentSubmissions(ctx, 10)
	if err != nil {
		s.logger.Error("failed to get submissionsActivities", "err", err)
		return nil, repository.ErrInternalServerError
	}

	return submissionsActivities, err
}

func (s *problemService) GetDailySubmissions(ctx context.Context, since int) ([]*models.DailySubmissions, error) {
	var dailySubmissions []*models.DailySubmissions
	var err error

	if since == 1 {
		dailySubmissions, err = s.submissionRepository.GetDailySubmissionsHours(ctx, 24)
	} else {
		dailySubmissions, err = s.submissionRepository.GetDailySubmissions(ctx, since)
	}

	if err != nil {
		s.logger.Error("failed to get daily submissions", "since", since, "err", err)
		return nil, repository.ErrInternalServerError
	}

	s.logger.Debug("daily submissions", "submissions", dailySubmissions)

	return dailySubmissions, nil
}

func (s *problemService) CreateNote(ctx context.Context, note *models.ProblemNote) error {
	note, err := s.problemNoteRepository.Create(ctx, note)
	if err != nil {
		s.logger.Error("repository: create note", "err", err, "note", note)
	}

	return err
}

func (s *problemService) DeleteNote(ctx context.Context, note *models.ProblemNote) error {
	err := s.problemNoteRepository.Delete(ctx, note.ID)
	if err != nil {
		s.logger.Error("repository: delete note", "err", err)
	}

	return err
}

func (s *problemService) GetNote(ctx context.Context, user *models.User, problem *models.Problem) (*models.ProblemNote, error) {
	problemNote, err := s.problemNoteRepository.Get(ctx, user.ID, problem.ID)
	if err != nil {
		s.logger.Error("repository: get problem note", "err", err)
		return nil, err
	}

	return problemNote, nil
}

func (s *problemService) UpdateNote(ctx context.Context, note *models.ProblemNote, markdown string) error {
	err := s.problemNoteRepository.UpdateMarkdown(ctx, note, markdown)
	if err != nil {
		s.logger.Error("repository: update note", "err", err)
	}

	return err
}

func (s *problemService) Search(ctx context.Context, search string) ([]*models.Problem, error) {
	s.logger.Debug("searching for problem", "search", search)

	problems, err := s.problemRepository.Search(
		ctx,
		search,
	)

	if err != nil {
		s.logger.Error("repository: search problem", "err", err)
		return nil, err
	}

	return problems, nil
}
