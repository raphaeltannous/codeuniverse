package judger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"git.riyt.dev/codeuniverse/internal/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var supportedLanguages = map[string]string{
	"c++":        "gcc:15.2-trixie",
	"golang":     "golang:1.25.4-alpine",
	"python3":    "python:3.13.7-alpine",
	"javascript": "node:22.21.1-alpine",
	"typescript": "node:22.21.1-alpine",
}

var (
	ErrLanguageNotSupported = errors.New("language not supported")
)

type Judge struct {
	Cli    *client.Client
	logger *slog.Logger
}

func NewJudge() (*Judge, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Judge{
		Cli:    cli,
		logger: slog.Default().With("package", "judge"),
	}, nil
}

func (judge *Judge) InitializeContainers(ctx context.Context) error {
	for _, wantedImage := range supportedLanguages {
		if err := judge.pullImageIfNotExists(ctx, wantedImage); err != nil {
			return err
		}
	}

	return nil
}

func (judge *Judge) pullImageIfNotExists(ctx context.Context, wantedImage string) error {
	imageInfo, err := judge.Cli.ImageInspect(ctx, wantedImage)
	if err != nil {
		return judge.pullImage(ctx, wantedImage)
	}

	containerLogger := judge.logger.WithGroup("containerInfo")
	containerLogger.Info(
		"already pulled.",
		"tag", wantedImage,
		"id", imageInfo.ID,
		"size", imageInfo.Size,
	)

	return nil
}

func (judge *Judge) pullImage(ctx context.Context, wantedImage string) error {
	judge.logger.Info("pulling image", "image", wantedImage)

	reader, err := judge.Cli.ImagePull(ctx, wantedImage, image.PullOptions{})
	if err != nil {
		judge.logger.Error("failed to pull image", "err", err)
		return err
	}
	defer reader.Close()

	_, err = io.Copy(logger.LoggerWriter, reader)
	if err != nil {
		judge.logger.Error("failed reading image stream", "err", err)
		return err
	}

	judge.logger.Info("image pulled", "image", wantedImage)

	return err
}

func (judge *Judge) Run(ctx context.Context, problemSlug, languageSlug, code string) error {
	if _, ok := supportedLanguages[languageSlug]; !ok {
		return ErrLanguageNotSupported
	}

	runWorkspace, err := os.MkdirTemp("", "run-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(runWorkspace)
	judge.logger.Debug("new runWorkspace", "workspace", runWorkspace)

	// test file
	mainTestFile := filepath.Join("problems-data", problemSlug, "go", "main_test.go")
	dstPath := filepath.Join(runWorkspace, "main_test.go")

	data, err := os.ReadFile(mainTestFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return err
	}

	// user file
	mainPath := filepath.Join(runWorkspace, "main.go")
	if err := os.WriteFile(mainPath, []byte(code), 0644); err != nil {
		return err
	}

	// go mod file

	modeContent := "module codeUniverse\n\ngo 1.25.4\n"
	modPath := filepath.Join(runWorkspace, "go.mod")
	if err := os.WriteFile(modPath, []byte(modeContent), 0644); err != nil {
		return err
	}

	resp, err := judge.Cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           supportedLanguages[languageSlug],
			Cmd:             []string{"go", "test", "-v", "."},
			NetworkDisabled: true,
			WorkingDir:      "/app",
		},
		&container.HostConfig{
			NetworkMode:    "none",
			ReadonlyRootfs: false,
			AutoRemove:     true,

			Resources: container.Resources{
				Memory:     256 * 1024 * 1024,
				MemorySwap: 256 * 1024 * 1024,
				CPUQuota:   50000,
				CPUPeriod:  100000,
			},

			Binds: []string{runWorkspace + ":/app:rw"},
		},
		nil,
		nil,
		"",
	)

	if err := judge.Cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := judge.Cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case status := <-statusCh:
		// status.StatusCode has the container exit code
		fmt.Println("Container exited with code:", status.StatusCode)
	case err := <-errCh:
		return fmt.Errorf("container wait error: %w", err)
	case <-ctx.Done():
		_ = judge.Cli.ContainerKill(context.Background(), resp.ID, "SIGKILL")
		return fmt.Errorf("execution timed out")
	}

	out, err := judge.Cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	judge.logger.Debug("achieved final destiny")
	return nil
}
