package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type javaJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (j *javaJudge) Run(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error) {
	return run(
		ctx,
		j.cli,
		problemCode,
		problemTestcases,
		&runConfig{
			cmd: []string{
				"sh", "-c",
				"javac -cp '.:gson.jar' *.java && java -cp '.:gson.jar' BackendChecker",
			},
		},
		j.logger,
	)
}

func newJavaJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &javaJudge{
		cli: cli,

		logger: logger,
	}
}
