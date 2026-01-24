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
	lessonsDir string
)

var (
	ErrNotFound = errors.New("static file not found")
)

func init() {
	// TODO: refactor
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
			log.Fatalf("failed to create directory %s: %s", courseDir, err)
		}
	}

	defaultCourseThumbnailDst := filepath.Join(courseDir, "default.jpg")
	slog.Info("Writing course-default.jpg", "destination", defaultCourseThumbnailDst)
	if err := os.WriteFile(defaultCourseThumbnailDst, defaultCourseThumbnail, 0644); err != nil {
		log.Fatalf("failed to write course-default.jpg to destination: %s", err)
	}

	lessonsDir = filepath.Join(staticDir, "lessons")
	if _, err := os.Stat(lessonsDir); errors.Is(err, os.ErrNotExist) {
		slog.Info("./static/lessons dir does not exists", "path", lessonsDir)
		if err := os.Mkdir(lessonsDir, 0770); err != nil {
			log.Fatalf("failed to create directory %s: %s", lessonsDir, err)
		}
	}
}

type StaticService interface {
	GetAvatar(ctx context.Context, filename string) (string, error)
	GetCourseThumbnail(ctx context.Context, filename string) (string, error)
	GetLessonVideo(ctx context.Context, filename string) (string, error)

	SaveAvatar(ctx context.Context, avatarSrc io.Reader, ext string) (string, error)
	DeleteAvatar(ctx context.Context, avatarPath string) error

	SaveCourseThumbnail(ctx context.Context, thumbnailSrc io.Reader, ext string) (string, error)
	DeleteCourseThumbnail(ctx context.Context, thumbnailPath string) error

	SaveLessonVideo(ctx context.Context, videoSrc io.Reader, ext string) (string, error)
	DeleteLessonVideo(ctx context.Context, videoPath string) error
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

// GetLessonVideo returns the path to be served using http.ServeFile.
func (s *staticService) GetLessonVideo(ctx context.Context, filename string) (string, error) {
	videoPath := filepath.Join(lessonsDir, filename)

	if _, err := os.Stat(videoPath); errors.Is(err, os.ErrNotExist) {
		s.logger.Debug("lesson video not found", "filename", filename, "videoPath", videoPath, "err", err)
		return "", ErrNotFound
	}

	s.logger.Debug("lesson found", "filename", filename, "filePath", videoPath)
	return videoPath, nil
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

func (s *staticService) SaveCourseThumbnail(ctx context.Context, thumbnailSrc io.Reader, ext string) (string, error) {
	uniqueFilename := uuid.New().String() + ext

	thumbnailPath := filepath.Join(courseDir, uniqueFilename)
	thumbnailDst, err := os.Create(thumbnailPath)
	if err != nil {
		s.logger.Error("failed to create thumbnailDst", "thumbnailPath", thumbnailPath, "err", err)
		return "", nil
	}
	defer thumbnailDst.Close()

	_, err = io.Copy(thumbnailDst, thumbnailSrc)

	return uniqueFilename, err
}

func (s *staticService) DeleteCourseThumbnail(ctx context.Context, thumbnailPath string) error {
	if thumbnailPath == "default.jpg" || thumbnailPath == "" {
		return nil
	}

	thumbnailPath = filepath.Join(courseDir, thumbnailPath)
	err := os.Remove(thumbnailPath)
	if err != nil {
		s.logger.Error("failed to delete thumbnailPath", "thumbnailPath", thumbnailPath)
		return err
	}

	return nil
}

func (s *staticService) DeleteLessonVideo(ctx context.Context, videoPath string) error {
	if videoPath == "default.mp4" || videoPath == "" {
		return nil
	}

	videoPath = filepath.Join(lessonsDir, videoPath)
	err := os.Remove(videoPath)
	if err != nil {
		s.logger.Error("failed to delete videoPath", "videoPath", videoPath)
		return err
	}

	return nil
}

func (s *staticService) SaveLessonVideo(ctx context.Context, videoSrc io.Reader, ext string) (string, error) {
	uniqueFilename := uuid.New().String() + ext

	videoPath := filepath.Join(lessonsDir, uniqueFilename)
	videoDst, err := os.Create(videoPath)
	if err != nil {
		s.logger.Error("failed to create videoDst", "videoPath", videoPath, "err", err)
		return "", nil
	}
	defer videoDst.Close()

	_, err = io.Copy(videoDst, videoSrc)

	return uniqueFilename, err
}
