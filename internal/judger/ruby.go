package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type rubyJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (r *rubyJudge) Run(
	ctx context.Context,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.RunResult, error) {
	panic("unimplemented")
}

func (r *rubyJudge) Submit(
	ctx context.Context,
	submission *models.Submission,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newRubyJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &rubyJudge{
		cli: cli,

		logger: logger,
	}
}
