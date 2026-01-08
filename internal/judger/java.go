package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type javaJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (j *javaJudge) Run(ctx context.Context, run *models.Run, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	panic("unimplemented")
}

func (j *javaJudge) Submit(ctx context.Context, submission *models.Submission, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newJavaJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &javaJudge{
		cli: cli,

		logger: logger,
	}
}
