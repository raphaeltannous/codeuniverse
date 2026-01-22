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

func (j *javascriptJudge) Run(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	return run(
		ctx,
		j.cli,
		problemCode,
		problemTestcases,
		&runConfig{
			cmd: []string{"node", "main.js"},
		},
		j.logger,
	)
}

func newJavascriptJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &javascriptJudge{
		cli: cli,

		logger: logger,
	}
}
