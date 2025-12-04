package judger

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type python3Judge struct {
	cli *client.Client

	logger *slog.Logger
}

func (p *python3Judge) Run(ctx context.Context, run *models.Run, problemSlug string) error {
	panic("unimplemented")
}

func (p *python3Judge) Submit(ctx context.Context, submission *models.Submission, problemSlug string) error {
	panic("unimplemented")
}

func newPython3Judge(cli *client.Client) languageJudge {
	return &python3Judge{
		cli: cli,

		logger: slog.Default().With("package", "judge.pythonJudge"),
	}
}
