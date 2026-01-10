package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type typescriptJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (t *typescriptJudge) Run(
	ctx context.Context,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.RunResult, error) {
	panic("unimplemented")
}

func (t *typescriptJudge) Submit(
	ctx context.Context,
	submission *models.Submission,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newTypescriptJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &typescriptJudge{
		cli: cli,

		logger: logger,
	}
}
