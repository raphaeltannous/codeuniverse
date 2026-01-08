package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type golangJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (g *golangJudge) Run(ctx context.Context, run *models.Run, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	panic("unimplemented")
}

func (g *golangJudge) Submit(ctx context.Context, submission *models.Submission, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newGolangJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &golangJudge{
		cli: cli,

		logger: logger,
	}
}
