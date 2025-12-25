package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type LessonRepository interface {
	Create(ctx context.Context, courseId uuid.UUID, lesson *models.Lesson) (*models.Lesson, error)
	Delete(ctx context.Context, lessonId uuid.UUID) error

	Get(ctx context.Context, lessonId uuid.UUID) (*models.Lesson, error)
	GetAllByCourse(ctx context.Context, courseId uuid.UUID) ([]*models.Lesson, error)

	UpdateTitle(ctx context.Context, lessonId uuid.UUID, title string) error
	UpdateDescription(ctx context.Context, lessonId uuid.UUID, description string) error
	UpdateVideoURL(ctx context.Context, lessonId uuid.UUID, url string) error
	UpdateDurationSeconds(ctx context.Context, lessonId uuid.UUID, duration int) error
	UpdateLessonNumber(ctx context.Context, lessonId uuid.UUID, lessonNumber int) error
}

var (
	ErrLessonNotFound              = errors.New("repository: lesson not found")
	ErrLessonNumberAlreadyAssigned = errors.New("repository: lesson number already assigned")
)
