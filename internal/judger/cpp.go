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

func (j *cppJudge) Run(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	return run(
		ctx,
		j.cli,
		problemCode,
		problemTestcases,
		&runConfig{
			cmd: []string{
				"sh", "-c",
				"g++ -std=c++23 -o program main.cpp && ./program",
			},
		},
		j.logger,
	)
}

func newCPPJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &cppJudge{
		cli: cli,

		logger: logger,
	}
}
