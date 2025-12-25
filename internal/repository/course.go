package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) (*models.Course, error)
	Delete(ctx context.Context, courseId uuid.UUID) error

	Get(ctx context.Context, courseId uuid.UUID) (*models.Course, error)

	UpdateTitle(ctx context.Context, courseId uuid.UUID, title string) error
	UpdateDescription(ctx context.Context, courseId uuid.UUID, description string) error
	UpdateThumbnailURL(ctx context.Context, courseId uuid.UUID, url string) error
}

var (
	ErrCourseNotFound = errors.New("repository: course not found")
)
