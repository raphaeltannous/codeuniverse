package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type cppJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (c *cppJudge) Run(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	panic("unimplemented")
}

func (c *cppJudge) Submit(ctx context.Context, submission *models.Submission, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newCPPJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &cppJudge{
		cli: cli,

		logger: logger,
	}
}
