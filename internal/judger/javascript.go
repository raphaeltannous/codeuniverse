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
	panic("unimplemented")
}

func newJavascriptJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &javascriptJudge{
		cli: cli,

		logger: logger,
	}
}
