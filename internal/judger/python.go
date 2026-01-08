package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type pythonJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (p *pythonJudge) Run(ctx context.Context, run *models.Run, problem *models.Problem, problemCode *models.ProblemCode) (*models.RunResult, error) {
	panic("unimplemented")
}

func (p *pythonJudge) Submit(ctx context.Context, submission *models.Submission, problem *models.Problem, problemCode *models.ProblemCode) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newPythonJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &pythonJudge{
		cli: cli,

		logger: logger,
	}
}
