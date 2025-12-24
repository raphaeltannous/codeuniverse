package services

import (
	"context"
	_ "embed"
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

//go:embed assets/default.png
var defaultUserProfilePicture []byte

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
		log.Fatalf("failed to get current working directory: %s", err)
	}
	staticDir = filepath.Join(currentDir, "static")
	if _, err := os.Stat(staticDir); errors.Is(err, os.ErrNotExist) {
		slog.Info("./static dir does not exists", "path", staticDir)
		if err := os.Mkdir(staticDir, 0770); err != nil {
			log.Fatalf("failed to create directory %s: %s", staticDir, err)
		}
	}

	avatarsDir = filepath.Join(staticDir, "avatars")
	if _, err := os.Stat(avatarsDir); errors.Is(err, os.ErrNotExist) {
		slog.Info("./static dir does not exists", "path", avatarsDir)
		if err := os.Mkdir(avatarsDir, 0770); err != nil {
			log.Fatalf("failed to create directory %s: %s", avatarsDir, err)
		}
	}

	slog.Info("writing default.png", "dstDir", avatarsDir)
	defaultAvatarDst := filepath.Join(avatarsDir, "default.png")

	if err := os.WriteFile(defaultAvatarDst, defaultUserProfilePicture, 0644); err != nil {
		log.Fatalf("failed to write default.png to destination: %s", err)
	}
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
