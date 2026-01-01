package repository

import (
	"context"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type CourseProgressRepository interface {
	Save(ctx context.Context, courseProgress *models.CourseLessonProgress) error

	Get(ctx context.Context, courseId, userId uuid.UUID) ([]*models.CourseLessonProgress, error)
}

var (
	ErrCourseProgressNotFound = errors.New("repository: course progress not found")
)
