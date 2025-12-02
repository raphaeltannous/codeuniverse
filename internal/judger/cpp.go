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

func newCPPJudge(cli *client.Client) languageJudge {
	return &cppJudge{
		cli: cli,

		logger: slog.Default().With("package", "judge.cppJudge"),
	}
}

func (c *cppJudge) Run(ctx context.Context, run *models.Run, problemSlug string) error {
	panic("unimplemented")
}

func (c *cppJudge) Submit(ctx context.Context, problemSlug string, code string) error {
	panic("unimplemented")
}
