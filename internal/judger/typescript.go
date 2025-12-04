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

func (t *typescriptJudge) Run(ctx context.Context, run *models.Run, problemSlug string) error {
	panic("unimplemented")
}

func (t *typescriptJudge) Submit(ctx context.Context, submission *models.Submission, problemSlug string) error {
	panic("unimplemented")
}

func newTypescriptJudge(cli *client.Client) languageJudge {
	return &typescriptJudge{
		cli: cli,

		logger: slog.Default().With("package", "judge.typescriptJudge"),
	}
}
