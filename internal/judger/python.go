package judger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type pythonJudge struct {
	cli *client.Client

	logger *slog.Logger
}

func (p *pythonJudge) Run(
	ctx context.Context,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.RunResult, error) {
	language := problemCode.Language

	runWorkspace, err := os.MkdirTemp("", fmt.Sprintf("run-%s-*", language.Slug()))
	if err != nil {
		p.logger.Error("failed to create runWorkspace", "runWorkspace", runWorkspace, "err", err)
		return &models.RunResult{Status: models.StatusInternalServerError}, err
	}
	defer os.RemoveAll(runWorkspace)
	p.logger.Debug("Created workspace", "runWorkspace", runWorkspace)
	p.logger.Debug("Setting up workspace.")

	if result, err := writeToWorkspace(runWorkspace, language.CheckerFilename(), problemCode.Checker, p.logger); err != nil {
		return result, err
	}

	if result, err := writeToWorkspace(runWorkspace, language.DriverFilename(), problemCode.Driver, p.logger); err != nil {
		return result, err
	}

	if result, err := writeToWorkspace(runWorkspace, language.CodeSnippetFilename(), problemCode.CodeSnippet, p.logger); err != nil {
		return result, err
	}

	p.logger.Debug("Copying runtime files", "files", problemCode.RuntimeFiles)
	if result, err := copyRuntimeFilesToWorkspace(runWorkspace, problemCode.RuntimeFiles); err != nil {
		p.logger.Error("Failed to copy runtime files", "runtimeFiles", problemCode.RuntimeFiles, "err", err)
		return result, err
	}

	if result, err := writeTestcasesToWorkspace(runWorkspace, problemTestcases); err != nil {
		return result, err
	}
	p.logger.Debug("Finished setting up runWorkspace.")

	resp, err := p.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           supportedLanguages[language].containerImage,
			Cmd:             []string{"python", language.BackendCheckerFilename()},
			NetworkDisabled: true,
			WorkingDir:      "/codeuniverse/run",
		},
		&container.HostConfig{
			Binds:         []string{runWorkspace + ":/codeuniverse/run:rw"},
			AutoRemove:    false,
			ReadonlyPaths: []string{},
			Resources: container.Resources{
				Memory: 256 * 1024 * 1024,
			},
		},
		nil,
		nil,
		"",
	)
	defer p.cli.ContainerRemove(context.WithoutCancel(ctx), resp.ID, container.RemoveOptions{Force: true})
	p.logger.Debug("Container created.", "container", resp)

	if err := p.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		p.logger.Error("failed to start container", "resp", resp, "err", err)
		return &models.RunResult{Status: models.StatusInternalServerError}, err
	}
	p.logger.Debug("Container started.")

	var peakMemory uint64
	statsCtx, cancelStats := context.WithCancel(ctx)
	defer cancelStats()

	go watchContainerMemory(statsCtx, p.cli, resp.ID, &peakMemory)

	statusChannel, errChannel := p.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case status := <-statusChannel:
		cancelStats()

		runResults := new(models.RunResult)
		switch status.StatusCode {
		case 0:
			jsonResults, err := os.Open(filepath.Join(runWorkspace, "results.json"))
			if err != nil {
				return &models.RunResult{Status: models.StatusInternalServerError}, err
			}

			results := new(models.Results)
			decoder := json.NewDecoder(jsonResults)
			if err := decoder.Decode(results); err != nil {
				return &models.RunResult{Status: models.StatusInternalServerError}, err
			}

			runResults = models.NewRunResult(results)
		case 1:
			runResults.Status = models.StatusCompileError
		}

		runResults.MemoryUsage = float64(peakMemory) / (1024 * 1024)
		stdOutErrSrc, err := p.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			p.logger.Error("failed to get stdOut and stdErr container logs", "err", err)
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}

		var stdOut strings.Builder
		var stdErr strings.Builder
		if _, err := stdcopy.StdCopy(&stdOut, &stdErr, stdOutErrSrc); err != nil {
			p.logger.Error("failed to write stdOut and stdErr container logs", "err", err)
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}

		runResults.StdOut = stdOut.String()
		runResults.StdErr = stdErr.String()

		return runResults, nil
	case err := <-errChannel:
		cancelStats()
		p.logger.Error("failed to wait for container", "err", err)

		return &models.RunResult{Status: models.StatusInternalServerError}, err
	case <-ctx.Done():
		p.logger.Debug("wait finished with ctx timeout")
		if err := p.cli.ContainerStop(context.WithoutCancel(ctx), resp.ID, container.StopOptions{}); err != nil {
			p.logger.Error("failed to stop container after ctx timeout", "err", err)
		}

		return &models.RunResult{Status: models.StatusTimeLimitExceeded}, nil
	}
}

func (p *pythonJudge) Submit(
	ctx context.Context,
	submission *models.Submission,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.SubmissionResult, error) {
	panic("unimplemented")
}

func newPythonJudge(cli *client.Client, logger *slog.Logger) languageJudger {
	return &pythonJudge{
		cli: cli,

		logger: logger,
	}
}
