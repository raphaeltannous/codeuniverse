package services

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error)
	DeleteCourse(ctx context.Context, course *models.Course) error

	GetCourseById(ctx context.Context, courseId uuid.UUID) (*models.Course, error)
	GetAllCourses(ctx context.Context) ([]*models.Course, error)

	UpdateCourse(ctx context.Context, courseId uuid.UUID, courseUpdatePatch map[string]any) (*models.Course, error)

	CreateLesson(ctx context.Context, course *models.Course, lesson *models.Lesson) (*models.Lesson, error)
	DeleteLesson(ctx context.Context, lesson *models.Lesson) error

	GetCourseLessons(ctx context.Context, courseId uuid.UUID) ([]*models.Lesson, error)
	GetLesson(ctx context.Context, lessonId uuid.UUID) (*models.Lesson, error)

	UpdateLesson(ctx context.Context, lessonId uuid.UUID, lessonUpdatePach map[string]any) (*models.Lesson, error)
}

type courseService struct {
	courseRepository repository.CourseRepository
	lessonRepository repository.LessonRepository

	logger *slog.Logger
}

func (c *courseService) CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	panic("unimplemented")
}

func (c *courseService) CreateLesson(ctx context.Context, course *models.Course, lesson *models.Lesson) (*models.Lesson, error) {
	panic("unimplemented")
}

func (c *courseService) DeleteCourse(ctx context.Context, course *models.Course) error {
	panic("unimplemented")
}

func (c *courseService) DeleteLesson(ctx context.Context, lesson *models.Lesson) error {
	panic("unimplemented")
}

func (c *courseService) GetAllCourses(ctx context.Context) ([]*models.Course, error) {
	panic("unimplemented")
}

func (c *courseService) GetCourseById(ctx context.Context, courseId uuid.UUID) (*models.Course, error) {
	panic("unimplemented")
}

func (c *courseService) GetCourseLessons(ctx context.Context, courseId uuid.UUID) ([]*models.Lesson, error) {
	panic("unimplemented")
}

func (c *courseService) GetLesson(ctx context.Context, lessonId uuid.UUID) (*models.Lesson, error) {
	panic("unimplemented")
}

func (c *courseService) UpdateCourse(ctx context.Context, courseId uuid.UUID, courseUpdatePatch map[string]any) (*models.Course, error) {
	panic("unimplemented")
}

func (c *courseService) UpdateLesson(ctx context.Context, lessonId uuid.UUID, lessonUpdatePach map[string]any) (*models.Lesson, error) {
	panic("unimplemented")
}

func NewCourseService(
	courseRepository repository.CourseRepository,
	lessonRepository repository.LessonRepository,

	logger *slog.Logger,
) CourseService {
	return &courseService{
		courseRepository: courseRepository,
		lessonRepository: lessonRepository,

		logger: slog.Default().With("package", "courseService"),
	}
}
