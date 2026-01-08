package judger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/logger"
	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

var (
	ErrLanguageNotSupported = errors.New("language not supported")
)

type languageJudger interface {
	Run(ctx context.Context, run *models.Run, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.RunResult, error)
	Submit(ctx context.Context, submission *models.Submission, problem *models.Problem, problemCode *models.ProblemCode, problemTestcases []*models.ProblemTestcase) (*models.SubmissionResult, error)
}

type Judge struct {
	cli       *client.Client
	logger    *slog.Logger
	languages map[models.ProblemLanguage]*LanguageJudge
}

func NewJudge() (*Judge, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Judge{
		cli:    cli,
		logger: slog.Default().With("package", "judge.Judge"),
	}, nil
}

func (judge *Judge) Close() error {
	return judge.cli.Close()
}

func (judge *Judge) InitializeContainers(ctx context.Context) error {
	for language, languageJudge := range supportedLanguages {
		if newFunc, ok := supportedLanguagesNewFunc[language]; ok {
			languageJudge.judge = newFunc(judge.cli, slog.Default().With("package", language.Slug()+"Language"))
		} else {
			return fmt.Errorf("failed to find newFunc for language %s", language)
		}

		if err := judge.pullImageIfNotExists(ctx, languageJudge.containerImage); err != nil {
			return err
		}
	}

	return nil
}

func (judge *Judge) pullImageIfNotExists(ctx context.Context, wantedImage string) error {
	imageInfo, err := judge.cli.ImageInspect(ctx, wantedImage)
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

	reader, err := judge.cli.ImagePull(ctx, wantedImage, image.PullOptions{})
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

func (judge *Judge) Run(
	ctx context.Context,
	run *models.Run,
	problem *models.Problem,
	problemCode *models.ProblemCode,
	problemTestcases []*models.ProblemTestcase,
) (*models.RunResult, error) {
	result, err := supportedLanguages[problemCode.Language].judge.Run(
		ctx,
		run,
		problem,
		problemCode,
		problemTestcases,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (judge *Judge) Submit(
	ctx context.Context,
	submission *models.Submission,
	problem *models.Problem,
	problemCode *models.ProblemCode,
) error {
	return nil
}
