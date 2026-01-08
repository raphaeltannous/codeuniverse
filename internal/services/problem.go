package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

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
	Delete(ctx context.Context, problem *models.Problem) error

	GetById(ctx context.Context, uuidString string) (*models.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*models.Problem, error)

	GetProblems(ctx context.Context, getParams *repository.GetProblemsParams) ([]*models.Problem, int, error)

	GetCount(ctx context.Context, difficulty models.ProblemDifficulty) (int, error)

	UpdateProblem(ctx context.Context, problem *models.Problem, problemUpdatePatch map[string]any) error

	Submit(ctx context.Context, user *models.User, problem *models.Problem, languageSlug, code string, handlerChannel chan string) error
	Run(ctx context.Context, user *models.User, problem *models.Problem, language models.ProblemLanguage, code string, handlerChannel chan string) error

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

	GetHint(ctx context.Context, id uuid.UUID) (*models.ProblemHint, error)
	GetProblemHints(ctx context.Context, problem *models.Problem) ([]*models.ProblemHint, error)
	CreateHint(ctx context.Context, problem *models.Problem, hint *models.ProblemHint) error
	UpdateHint(ctx context.Context, hint *models.ProblemHint, content string) error
	DeleteHint(ctx context.Context, hint *models.ProblemHint) error

	GetProblemCodes(ctx context.Context, problem *models.Problem) ([]*models.ProblemCode, error)
	SaveProblemCode(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode) error

	GetProblemTestcases(ctx context.Context, problem *models.Problem) ([]*models.ProblemTestcase, error)
	AddProblemTestcase(ctx context.Context, problem *models.Problem, testcase *models.ProblemTestcase) error
	UpdateProblemTestcase(ctx context.Context, problem *models.Problem, testcaseId int, testcase *models.ProblemTestcase) error
	DeleteProblemTestcase(ctx context.Context, problem *models.Problem, testcaseId int) error

	GetProblemCodeConfig(ctx context.Context, problem *models.Problem) (*models.ProblemCodeConfig, error)
	UpdateProblemCodeConfig(ctx context.Context, problem *models.Problem, config *models.ProblemCodeConfig) error
}

type problemService struct {
	problemRepository     repository.ProblemRepository
	problemNoteRepository repository.ProblemNoteRepository
	runRepository         repository.RunRepository
	submissionRepository  repository.SubmissionRepository
	problemHintRepository repository.ProblemHintRepository
	problemCodeRepository repository.ProblemCodeRepository

	judge  judger.Judge
	logger *slog.Logger
}

func (s *problemService) AddProblemTestcase(
	ctx context.Context,
	problem *models.Problem,
	testcase *models.ProblemTestcase,
) error {
	problemTestcases, err := s.GetProblemTestcases(ctx, problem)
	if err != nil {
		return err
	}

	maxTestcaseId := 0
	for _, problemTestcase := range problemTestcases {
		maxTestcaseId = max(maxTestcaseId, problemTestcase.Id)
	}

	testcase.Id = maxTestcaseId + 1
	problemTestcases = append(problemTestcases, testcase)
	err = s.problemCodeRepository.SaveTestcases(
		ctx,
		problem,
		problemTestcases,
	)
	if err != nil {
		s.logger.Error("failed to add testcase", "problem", problem, "testcase", testcase, "err", err)
		return err
	}

	return nil
}

func (s *problemService) DeleteProblemTestcase(
	ctx context.Context,
	problem *models.Problem,
	testcaseId int,
) error {
	problemTestcases, err := s.GetProblemTestcases(ctx, problem)
	if err != nil {
		return err
	}

	problemTestcases = slices.DeleteFunc(problemTestcases, func(a *models.ProblemTestcase) bool {
		return a.Id == testcaseId
	})

	err = s.problemCodeRepository.SaveTestcases(
		ctx,
		problem,
		problemTestcases,
	)
	if err != nil {
		s.logger.Error("failed to delete problem code testcase", "problem", problem, "testcaseId", testcaseId, "err", err)
		return err
	}

	return nil
}

func (s *problemService) GetProblemCodeConfig(ctx context.Context, problem *models.Problem) (*models.ProblemCodeConfig, error) {
	config, err := s.problemCodeRepository.GetProblemCodeConfig(
		ctx,
		problem,
	)
	if err != nil {
		s.logger.Error("failed to get problem code config", "problem", problem, "err", err)
		return nil, err
	}

	return config, nil
}

func (s *problemService) GetProblemTestcases(ctx context.Context, problem *models.Problem) ([]*models.ProblemTestcase, error) {
	problemTestcases, err := s.problemCodeRepository.GetTestcases(
		ctx,
		problem,
	)
	if err != nil {
		s.logger.Error("failed to get problem code testcases", "problem", problem, "err", err)
		return nil, nil
	}

	return problemTestcases, nil
}

func (s *problemService) UpdateProblemCodeConfig(ctx context.Context, problem *models.Problem, config *models.ProblemCodeConfig) error {
	err := s.problemCodeRepository.SaveProblemCodeConfig(
		ctx,
		problem,
		config,
	)
	if err != nil {
		s.logger.Error("failed to save new problem code config", "problem", problem, "config", config, "err", err)
		return err
	}

	return nil
}

func (s *problemService) UpdateProblemTestcase(
	ctx context.Context,
	problem *models.Problem,
	testcaseId int,
	testcase *models.ProblemTestcase,
) error {
	problemTestcases, err := s.GetProblemTestcases(ctx, problem)
	if err != nil {
		return err
	}

	for i, problemTestcase := range problemTestcases {
		if problemTestcase.Id == testcaseId {
			testcase.Id = testcaseId
			problemTestcases[i] = testcase
			break
		}
	}

	err = s.problemCodeRepository.SaveTestcases(
		ctx,
		problem,
		problemTestcases,
	)
	if err != nil {
		s.logger.Error("failed to update problem code testcase", "problem", problem, "testcaseId", testcaseId, "testcase", testcase, "err", err)
		return err
	}

	return nil
}

func NewProblemService(
	problemRepository repository.ProblemRepository,
	problemNoteRepository repository.ProblemNoteRepository,
	runRepository repository.RunRepository,
	submissionRepository repository.SubmissionRepository,
	problemHintRepository repository.ProblemHintRepository,
	problemCodeRepository repository.ProblemCodeRepository,

	judge judger.Judge,
) ProblemService {
	return &problemService{
		problemRepository:     problemRepository,
		problemNoteRepository: problemNoteRepository,
		runRepository:         runRepository,
		submissionRepository:  submissionRepository,
		problemHintRepository: problemHintRepository,
		problemCodeRepository: problemCodeRepository,

		judge:  judge,
		logger: slog.Default().With("package", "problemsService"),
	}
}

func (s *problemService) Create(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	problem, err := s.problemRepository.Create(ctx, problem)
	if err != nil {
		s.logger.Error("failed to create problem", "problem", problem, "err", err)
		return nil, repository.ErrInternalServerError
	}

	return problem, nil
}

func (s *problemService) Delete(ctx context.Context, problem *models.Problem) error {
	err := s.problemRepository.Delete(
		ctx,
		problem.ID,
	)
	if err != nil {
		s.logger.Error("failed to delete problem", "problem", problem, "err", err)
		return err
	}

	if problem.ID.String() != "" {
		return s.problemCodeRepository.DeleteProblem(ctx, problem)
	}

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

	return problem, nil
}

func (s *problemService) GetProblems(ctx context.Context, getParams *repository.GetProblemsParams) ([]*models.Problem, int, error) {
	problems, total, err := s.problemRepository.GetProblems(
		ctx,
		getParams,
	)
	if err != nil {
		s.logger.Error("failed to get problems", "getParams", getParams, "err", err)
		return nil, 0, repository.ErrInternalServerError
	}
	return problems, total, nil
}

func (s *problemService) GetCount(ctx context.Context, difficulty models.ProblemDifficulty) (int, error) {
	count, err := s.problemRepository.GetCountByDifficulty(ctx, difficulty)
	if err != nil {
		s.logger.Error("failed to get easy count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *problemService) UpdateProblem(
	ctx context.Context,
	problem *models.Problem,
	updateProblemPatch map[string]any,
) error {
	if err := updateProblemFromPatch(ctx, problem, updateProblemPatch, "title", s.problemRepository.UpdateTitle); err != nil {
		s.logger.Error("failed to update title", "problem", problem, "updateProblemPatch", updateProblemPatch, "err", err)
		return err
	}

	if err := updateProblemFromPatch(ctx, problem, updateProblemPatch, "slug", s.problemRepository.UpdateSlug); err != nil {
		s.logger.Error("failed to update slug", "problem", problem, "updateProblemPatch", updateProblemPatch, "err", err)
		return err
	}

	if err := updateProblemFromPatch(ctx, problem, updateProblemPatch, "description", s.problemRepository.UpdateDescription); err != nil {
		s.logger.Error("failed to update description", "problem", problem, "updateProblemPatch", updateProblemPatch, "err", err)
		return err
	}

	if err := updateProblemFromPatch(ctx, problem, updateProblemPatch, "difficulty", s.problemRepository.UpdateDifficulty); err != nil {
		s.logger.Error("failed to update difficulty", "problem", problem, "updateProblemPatch", updateProblemPatch, "err", err)
		return err
	}

	if err := updateProblemFromPatch(ctx, problem, updateProblemPatch, "isPremium", s.problemRepository.UpdateIsPremium); err != nil {
		s.logger.Error("failed to update isPremium", "problem", problem, "updateProblemPatch", updateProblemPatch, "err", err)
		return err
	}

	if err := updateProblemFromPatch(ctx, problem, updateProblemPatch, "isPublic", s.problemRepository.UpdateIsPublic); err != nil {
		s.logger.Error("failed to update isPublic", "problem", problem, "updateProblemPatch", updateProblemPatch, "err", err)
		return err
	}

	return nil
}

func updateProblemFromPatch[T any](
	ctx context.Context,
	problem *models.Problem,
	patch map[string]any,
	key string,
	updateFunc func(ctx context.Context, id uuid.UUID, value T) error,
) error {
	if rawValue, ok := patch[key]; ok {
		value, ok := rawValue.(T)
		if !ok {
			return fmt.Errorf("invalid type for %s: %T", key, rawValue)
		}

		return updateFunc(ctx, problem.ID, value)
	}

	return nil
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
		problem,
		nil, // TODO
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

func (s *problemService) Run(
	ctx context.Context,
	user *models.User,
	problem *models.Problem,
	language models.ProblemLanguage,
	code string,
	handlerChannel chan string,
) error {
	run := models.NewRun(
		user.ID,
		problem.ID,
		language.Slug(),
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
		problem,
		nil, // TODO
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

func (s *problemService) CreateHint(
	ctx context.Context,
	problem *models.Problem,
	hint *models.ProblemHint,
) error {
	hint.ProblemId = problem.ID
	if err := s.problemHintRepository.Create(ctx, hint); err != nil {
		s.logger.Error("failed to create hint for problem", "problem", problem, "hint", hint, "err", err)
		return err
	}

	return nil
}

func (s *problemService) DeleteHint(
	ctx context.Context,
	hint *models.ProblemHint,
) error {
	if err := s.problemHintRepository.Delete(ctx, hint.ID); err != nil {
		s.logger.Error("failed to delete hint", "hint", hint, "err", err)
		return err
	}

	return nil
}

func (s *problemService) GetHint(
	ctx context.Context,
	id uuid.UUID,
) (*models.ProblemHint, error) {
	hint, err := s.problemHintRepository.Get(
		ctx,
		id,
	)
	if err != nil {
		s.logger.Error("failed to get hint", "id", id, "err", err)
		return nil, err
	}

	return hint, nil
}

func (s *problemService) GetProblemHints(
	ctx context.Context,
	problem *models.Problem,
) ([]*models.ProblemHint, error) {
	hints, err := s.problemHintRepository.GetHints(
		ctx,
		problem.ID,
	)
	if err != nil {
		s.logger.Error("failed to get hints for problem", "problem", problem, "err", err)
		return []*models.ProblemHint(nil), err
	}

	return hints, nil
}

func (s *problemService) UpdateHint(
	ctx context.Context,
	hint *models.ProblemHint,
	content string,
) error {
	if err := s.problemHintRepository.Update(ctx, hint.ID, content); err != nil {
		s.logger.Error("failed to update hint", "hint", hint, "content", content, "err", err)
		return err
	}

	return nil
}

func (s *problemService) GetProblemCodes(ctx context.Context, problem *models.Problem) ([]*models.ProblemCode, error) {
	problemCodes, err := s.problemCodeRepository.GetProblemCodes(ctx, problem)
	if err != nil {
		s.logger.Error("failed to get problemCodes", "problem", problem, "err", err)
		return nil, err
	}

	return problemCodes, nil
}

func (s *problemService) SaveProblemCode(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode) error {
	err := s.problemCodeRepository.SaveProblemCode(
		ctx,
		problem,
		problemCode,
	)
	if err != nil {
		s.logger.Error("failed to save problemCode", "problem", problem, "problemCode", problemCode, "err", err)
		return err
	}

	return nil
}
