package judger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type golangJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (g *golangJudge) Run(ctx context.Context, run *models.Run, problemSlug string) error {
	runWorkspace, err := os.MkdirTemp("", "run-go-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(runWorkspace)
	g.logger.Debug("creating run workspace", "runWorkspace", runWorkspace)

	l, ok := SupportedLanguages["golang"]
	if !ok {
		return errors.New("unsupported language")
	}
	l.copyRunToWorkspace(problemSlug, runWorkspace)

	userSubmittedCode := filepath.Join(runWorkspace, "main.go")
	if run.Code[:12] != "package main" {
		run.Code = "package main\n" + run.Code
	}
	if err := os.WriteFile(userSubmittedCode, []byte(run.Code), 0644); err != nil {
		return err
	}

	modFile := filepath.Join(runWorkspace, "go.mod")
	if err := os.WriteFile(modFile, g.getModContent(problemSlug), 0644); err != nil {
		return err
	}
	g.logger.Debug("finished settinp up workspace")

	resp, err := g.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           SupportedLanguages[run.Language].containerImage,
			Cmd:             []string{"go", "test", "-timeout", "7s", "."},
			NetworkDisabled: true,
			WorkingDir:      "/app",
		},
		&container.HostConfig{
			AutoRemove: true,
			Binds:      []string{runWorkspace + ":/app:rw"},
		},
		nil,
		nil,
		"",
	)

	if err := g.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
	g.logger.Debug("started container")

	statusCh, errCh := g.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case status := <-statusCh:
		g.logger.Warn("status", "status", status)
		if status.StatusCode == 0 {
			run.IsAccepted = true
			run.Status = "ACCEPTED"
		} else {
			run.IsAccepted = false
			run.Status = "FAILED"
		}
	case <-errCh:
		run.IsAccepted = false
		run.Status = "Internal Server Error. Please Contact Support."
	case <-ctx.Done():
		run.IsAccepted = false
		run.Status = "TIME EXCEED LIMIT"
	}
	g.logger.Debug("container is finished")

	out, err := g.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}

func (g *golangJudge) Submit(ctx context.Context, submission *models.Submission, problemSlug string) error {
	submitWorkspace, err := os.MkdirTemp("", "submit-go-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(submitWorkspace)
	g.logger.Debug("creating run workspace", "runWorkspace", submitWorkspace)

	l, ok := SupportedLanguages["golang"]
	if !ok {
		return errors.New("unsupported language")
	}
	l.copySubmitToWorkspace(problemSlug, submitWorkspace)

	userSubmittedCode := filepath.Join(submitWorkspace, "main.go")
	if submission.Code[:12] != "package main" {
		submission.Code = "package main\n" + submission.Code
	}
	if err := os.WriteFile(userSubmittedCode, []byte(submission.Code), 0644); err != nil {
		return err
	}

	modFile := filepath.Join(submitWorkspace, "go.mod")
	if err := os.WriteFile(modFile, g.getModContent(problemSlug), 0644); err != nil {
		return err
	}
	g.logger.Debug("finished settinp up workspace")

	resp, err := g.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           SupportedLanguages[submission.Language].containerImage,
			Cmd:             []string{"go", "test", "-timeout", "7s", "."},
			NetworkDisabled: true,
			WorkingDir:      "/app",
		},
		&container.HostConfig{
			AutoRemove: true,
			Binds:      []string{submitWorkspace + ":/app:rw"},
		},
		nil,
		nil,
		"",
	)

	if err := g.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
	g.logger.Debug("started container")

	statusCh, errCh := g.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case status := <-statusCh:
		g.logger.Warn("status", "status", status)
		if status.StatusCode == 0 {
			submission.IsAccepted = true
			submission.Status = "ACCEPTED"
		} else {
			submission.IsAccepted = false
			submission.Status = "FAILED"
		}
	case <-errCh:
		submission.IsAccepted = false
		submission.Status = "Internal Server Error. Please Contact Support."
	case <-ctx.Done():
		submission.IsAccepted = false
		submission.Status = "TIME EXCEED LIMIT"
	}
	g.logger.Debug("container is finished")

	out, err := g.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}

func (g *golangJudge) getModContent(moduleName string) []byte {
	return fmt.Appendf(nil, "module %s\n\ngo 1.25.4", moduleName)
}

func newGolangJudge(cli *client.Client) languageJudge {
	return &golangJudge{
		cli: cli,

		logger: slog.Default().With("package", "judge.golangJudge"),
	}
}
