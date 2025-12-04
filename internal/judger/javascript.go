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

func (j *javascriptJudge) Run(ctx context.Context, run *models.Run, problemSlug string) error {
	panic("unimplemented")
}

func (j *javascriptJudge) Submit(ctx context.Context, submission *models.Submission, problemSlug string) error {
	panic("unimplemented")
}

func newJavascriptJudge(cli *client.Client) languageJudge {
	return &javascriptJudge{
		cli: cli,

		logger: slog.Default().With("package", "judge.javascriptJudge"),
	}
}
