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

	GetBySlug(ctx context.Context, slug string) (*models.Course, error)
	GetAll(ctx context.Context) ([]*models.Course, error)
	GetAllPublished(ctx context.Context) ([]*models.Course, error)

	UpdateTitle(ctx context.Context, courseId uuid.UUID, title string) error
	UpdateSlug(ctx context.Context, courseId uuid.UUID, slug string) error
	UpdateDifficulty(ctx context.Context, courseId uuid.UUID, difficulty string) error
	UpdateDescription(ctx context.Context, courseId uuid.UUID, description string) error
	UpdateThumbnailURL(ctx context.Context, courseId uuid.UUID, url string) error
	UpdateIsPublished(ctx context.Context, courseId uuid.UUID, status bool) error
}

var (
	ErrCourseNotFound      = errors.New("repository: course not found")
	ErrCourseAlreadyExists = errors.New("repository: course already exists")
)
