package services

import (
	"context"
	_ "embed"
	"errors"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

//go:embed assets/avatar-default.png
var defaultUserProfilePicture []byte

//go:embed assets/course-default.jpg
var defaultCourseThumbnail []byte

var (
	staticDir  string
	avatarsDir string
	courseDir  string
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
		slog.Info("./static/avatars dir does not exists", "path", avatarsDir)
		if err := os.Mkdir(avatarsDir, 0770); err != nil {
			log.Fatalf("failed to create directory %s: %s", avatarsDir, err)
		}
	}

	defaultAvatarDst := filepath.Join(avatarsDir, "default.png")
	slog.Info("writing avatar-default.png", "destination", defaultAvatarDst)
	if err := os.WriteFile(defaultAvatarDst, defaultUserProfilePicture, 0644); err != nil {
		log.Fatalf("failed to write avatar-default.png to destination: %s", err)
	}

	courseDir = filepath.Join(staticDir, "courses")
	if _, err := os.Stat(courseDir); errors.Is(err, os.ErrNotExist) {
		slog.Info("./static/courses dir does not exists", "path", courseDir)
		if err := os.Mkdir(courseDir, 0770); err != nil {
			log.Fatalf("failed to create directory %s: %s", avatarsDir, err)
		}
	}

	defaultCourseThumbnailDst := filepath.Join(courseDir, "default.jpg")
	slog.Info("Writing course-default.jpg", "destination", defaultCourseThumbnailDst)
	if err := os.WriteFile(defaultCourseThumbnailDst, defaultCourseThumbnail, 0644); err != nil {
		log.Fatalf("failed to write course-default.jpg to destination: %s", err)
	}
}

type StaticService interface {
	GetAvatar(ctx context.Context, filename string) (string, error)
	GetCourseThumbnail(ctx context.Context, filename string) (string, error)

	SaveAvatar(ctx context.Context, avatarSrc io.Reader, ext string) (string, error)
	DeleteAvatar(ctx context.Context, avatarPath string) error
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

// GetCourseThumbnail returns the path to be served using http.ServeFile.
func (s *staticService) GetCourseThumbnail(ctx context.Context, filename string) (string, error) {
	filePath := filepath.Join(courseDir, filename)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		s.logger.Debug("course thumbnail not found", "filename", filename, "filePath", filePath, "err", err)
		return "", ErrNotFound
	}

	s.logger.Debug("thumbnail found", "filename", filename, "filePath", filePath)
	return filePath, nil
}

func (s *staticService) SaveAvatar(ctx context.Context, avatarSrc io.Reader, ext string) (string, error) {
	uniqueFilename := uuid.New().String() + ext

	avatarPath := filepath.Join(avatarsDir, uniqueFilename)
	avatarDst, err := os.Create(avatarPath)
	if err != nil {
		s.logger.Error("failed to create avatarDst", "avatarPath", avatarPath, "err", err)
		return "", nil
	}
	defer avatarDst.Close()

	_, err = io.Copy(avatarDst, avatarSrc)

	return uniqueFilename, err
}

func (s *staticService) DeleteAvatar(ctx context.Context, avatarPath string) error {
	if avatarPath == "default.png" || avatarPath == "" {
		return nil
	}

	avatarPath = filepath.Join(avatarsDir, avatarPath)
	err := os.Remove(avatarPath)
	if err != nil {
		s.logger.Error("failed to delete avatarPath", "avatarPath", avatarPath)
		return err
	}

	return nil
}
