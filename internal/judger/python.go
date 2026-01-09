package judger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
	run *models.Run,
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

	writeFunc := func(filename string, content string) (*models.RunResult, error) {
		filePath := filepath.Join(runWorkspace, filename)
		if err := os.WriteFile(filePath, []byte(content), 0640); err != nil {
			p.logger.Error("failed to write file", "filename", filename, "content", content, "err", err)
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}

		return nil, nil
	}

	if result, err := writeFunc(language.CheckerFilename(), problemCode.Checker); err != nil {
		return result, err
	}

	if result, err := writeFunc(language.DriverFilename(), problemCode.Driver); err != nil {
		return result, err
	}

	if result, err := writeFunc(language.CodeSnippetFilename(), problemCode.CodeSnippet); err != nil {
		return result, err
	}

	writeTestcasesFunc := func(testcases []*models.ProblemTestcase) (*models.RunResult, error) {
		testcasesPath := filepath.Join(runWorkspace, "testcases.json")
		testcasesFile, err := os.Create(testcasesPath)
		if err != nil {
			p.logger.Error("failed to create tescases file", "testcasesPath", testcasesPath, "err", err)
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}
		defer testcasesFile.Close()

		encoder := json.NewEncoder(testcasesFile)
		if err := encoder.Encode(testcases); err != nil {
			p.logger.Error("failed to write testcases to file", "testcasesPath", testcasesPath, "err", err)
			return &models.RunResult{Status: models.StatusInternalServerError}, err
		}

		return nil, nil
	}

	if result, err := writeTestcasesFunc(problemTestcases); err != nil {
		return result, err
	}
	p.logger.Debug("Finished setting up runWorkspace.")

	resp, err := p.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           supportedLanguages[language].containerImage,
			Cmd:             []string{"python", language.CheckerFilename()},
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

	go func() {
		stats, err := p.cli.ContainerStats(statsCtx, resp.ID, true)
		if err != nil {
			return
		}
		defer stats.Body.Close()

		decoder := json.NewDecoder(stats.Body)
		for {
			var v container.StatsResponse
			if err := decoder.Decode(&v); err != nil {
				return
			}
			if v.MemoryStats.Usage > peakMemory {
				peakMemory = v.MemoryStats.Usage
			}
		}
	}()

	statusChannel, errChannel := p.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case status := <-statusChannel:
		p.logger.Debug("container finised", "status", status)
		out, err := p.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			panic(err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		cancelStats()
	case err := <-errChannel:
		out, err := p.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			panic(err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		cancelStats()
		p.logger.Debug("container finished", "err", err)
	case <-ctx.Done():
		p.logger.Debug("wait finished with ctx timeout")
		if err := p.cli.ContainerStop(context.WithoutCancel(ctx), resp.ID, container.StopOptions{}); err != nil {
			p.logger.Error("failed to stop container after ctx timeout", "err", err)
		}
	}

	p.logger.Debug("Container stopped")
	p.logger.Debug("final peak memory", "usage", float64(peakMemory)/(1024*1024))

	return nil, nil
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
