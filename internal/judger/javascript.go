package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type javascriptJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (j *javascriptJudge) Run(ctx context.Context, run *models.Run, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	panic("unimplemented")
}

func (j *javascriptJudge) Submit(ctx context.Context, submission *models.Submission, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newJavascriptJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &javascriptJudge{
		cli: cli,

		logger: logger,
	}
}
