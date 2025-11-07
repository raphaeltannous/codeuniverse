package judger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

var supportedLanguages = map[string]string{
	"golang":     "golang:1.25.4-alpine",
	"python":     "python:3.13.7-alpine",
	"c++":        "gcc:15.2-trixie",
	"javascript": "node:22.21.1-alpine",
}

type Judge struct {
	Cli *client.Client
}

func NewJudge() (*Judge, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Judge{
		Cli: cli,
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

	fmt.Println(wantedImage, imageInfo.ID, imageInfo.Size)

	return nil
}

func (judge *Judge) pullImage(ctx context.Context, wantedImage string) error {
	reader, err := judge.Cli.ImagePull(ctx, wantedImage, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)

	return err
}
