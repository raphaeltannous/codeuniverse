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

func (j *pythonJudge) Run(
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
			cmd: []string{"python", problemCode.Language.BackendCheckerFilename()},
		},
		j.logger,
	)
}

func newPythonJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &pythonJudge{
		cli: cli,

		logger: logger,
	}
}
