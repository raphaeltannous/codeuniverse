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

func (j *golangJudge) Run(
	ctx context.Context,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.RunResult, error) {
	return run(
		ctx,
		j.cli,
		problemCode,
		problemTestcases,
		&runConfig{
			cmd: []string{
				"go", "run", ".",
			},
		},
		j.logger,
	)
}

func newGolangJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &golangJudge{
		cli: cli,

		logger: logger,
	}
}
