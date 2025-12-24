package services

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	staticDir  string
	avatarsDir string
)

var (
	ErrNotFound = errors.New("static file not found")
)

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory: %w", err)
	}
	staticDir = filepath.Join(currentDir, "static")

	avatarsDir = filepath.Join(staticDir, "avatars")
}

type StaticService interface {
	GetAvatar(ctx context.Context, filename string) (string, error)
}

type staticService struct {
	logger *slog.Logger
}

func NewStaticService() StaticService {
	return &staticService{
		logger: slog.Default().With("package", "staticService"),
	}
}

// GetAvatar returns the path to be served using http.ServeFile.
func (s *staticService) GetAvatar(ctx context.Context, filename string) (string, error) {
	filePath := filepath.Join(avatarsDir, filename)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		s.logger.Debug("avatar not found", "filename", filename, "err", err)
		return "", ErrNotFound
	}
	s.logger.Debug("avatar found", "filename", filename, "filePath", filePath)

	return filePath, nil
}
