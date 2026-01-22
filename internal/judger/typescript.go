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

func (j *typescriptJudge) Run(
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
			cmd: []string{"npx", "tsx", "main.ts"},
		},
		j.logger,
	)
}

func newTypescriptJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &typescriptJudge{
		cli: cli,

		logger: logger,
	}
}
