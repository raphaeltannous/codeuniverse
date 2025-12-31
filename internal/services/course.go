package services

import (
	"context"
	"errors"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/utils"
	"github.com/google/uuid"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error)
	DeleteCourse(ctx context.Context, course *models.Course) error

	GetCourseBySlug(ctx context.Context, slug string) (*models.Course, error)
	GetAllCourses(ctx context.Context) ([]*models.Course, error)

	UpdateCourse(ctx context.Context, course *models.Course, courseUpdatePatch map[string]any) error

	CreateLesson(ctx context.Context, course *models.Course, lesson *models.Lesson) (*models.Lesson, error)
	DeleteLesson(ctx context.Context, lesson *models.Lesson) error

	GetCourseLessons(ctx context.Context, course *models.Course) ([]*models.Lesson, error)
	GetLesson(ctx context.Context, lessonId uuid.UUID) (*models.Lesson, error)

	UpdateLesson(ctx context.Context, lesson *models.Lesson, lessonUpdatePach map[string]any) error
}

var (
	ErrInvalidPatch = errors.New("service: invalid patch")
)

type courseService struct {
	courseRepository repository.CourseRepository
	lessonRepository repository.LessonRepository

	logger *slog.Logger
}

func (c *courseService) GetAllCourses(ctx context.Context) ([]*models.Course, error) {
	courses, err := c.courseRepository.GetAll(ctx)
	if err != nil {
		c.logger.Error("failed to get all courses", "err", err)
		return nil, repository.ErrInternalServerError
	}

	for _, course := range courses {
		count, err := c.lessonRepository.GetLessonCountForCourse(ctx, course.ID)
		if err != nil {
			c.logger.Error("failed to get lesson count for course", "course", course, "err", err)
			return nil, repository.ErrInternalServerError
		}
		course.TotalLessons = count
	}

	return courses, nil
}

func (c *courseService) CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	course, err := c.courseRepository.Create(ctx, course)
	if err != nil {
		c.logger.Error("failed to create course", "course", course, "err", err)
		return nil, err
	}

	return course, nil
}

func (c *courseService) CreateLesson(ctx context.Context, course *models.Course, lesson *models.Lesson) (*models.Lesson, error) {
	lesson, err := c.lessonRepository.Create(
		ctx,
		course.ID,
		lesson,
	)
	if err != nil {
		c.logger.Error("failed to create lesson for course", "course", course, "lesson", lesson, "err", err)
		return nil, err
	}

	return lesson, nil
}

func (c *courseService) DeleteCourse(ctx context.Context, course *models.Course) error {
	err := c.courseRepository.Delete(ctx, course.ID)
	if err != nil {
		c.logger.Error("failed to delete course", "course", course)
		return err
	}

	return nil
}

func (c *courseService) DeleteLesson(ctx context.Context, lesson *models.Lesson) error {
	err := c.lessonRepository.Delete(ctx, lesson.ID)
	if err != nil {
		c.logger.Error("failed to delete lesson", "lesson", lesson)
		return err
	}

	return nil
}

func (c *courseService) GetCourseBySlug(ctx context.Context, slug string) (*models.Course, error) {
	course, err := c.courseRepository.GetBySlug(ctx, slug)
	if err != nil {
		c.logger.Error("failed to get course by slug", "slug", slug, "err", err)
		return nil, repository.ErrInternalServerError
	}

	return course, nil
}

func (c *courseService) GetCourseLessons(ctx context.Context, course *models.Course) ([]*models.Lesson, error) {
	lessons, err := c.lessonRepository.GetAllByCourse(ctx, course.ID)
	if err != nil {
		c.logger.Error("failed to get lessons for course", "course", course, "err", err)
		return nil, err
	}

	for _, lesson := range lessons {
		if lesson.VideoURL != "default.mp4" {
			lesson.VideoURL = utils.GenerateSignedUrl(lesson.ID.String(), 60*24)
		}
	}

	return lessons, nil
}

func (c *courseService) GetLesson(ctx context.Context, lessonId uuid.UUID) (*models.Lesson, error) {
	lesson, err := c.lessonRepository.Get(ctx, lessonId)
	if err != nil {
		c.logger.Error("failed to get lesson", "lessonId", lessonId, "err", err)
		return nil, err
	}

	return lesson, nil
}

func (c *courseService) UpdateCourse(ctx context.Context, course *models.Course, courseUpdatePatch map[string]any) error {
	if rawThumbnailUrl, ok := courseUpdatePatch["thumbnailUrl"]; ok {
		switch thumbnailUrl := rawThumbnailUrl.(type) {
		case string:
			err := c.courseRepository.UpdateThumbnailURL(ctx, course.ID, thumbnailUrl)
			if err != nil {
				c.logger.Error("failed to update thumbnail_url", "course", course, "newThumbnailUrl", thumbnailUrl, "err", err)
				return err
			}
		default:
			c.logger.Error("thumbnailUrl is not a string", "rawThumbnailUrl", rawThumbnailUrl)
			return ErrInvalidPatch
		}

		return nil
	}

	if rawIsPublished, ok := courseUpdatePatch["isPublished"]; ok {
		switch isPublished := rawIsPublished.(type) {
		case bool:
			err := c.courseRepository.UpdateIsPublished(ctx, course.ID, isPublished)
			if err != nil {
				c.logger.Error("failed to update is_published status", "course", course, "newStatus", isPublished, "err", err)
				return err
			}
		default:
			c.logger.Error("isPublished is not a bool", "rawIsPublished", rawIsPublished)
			return ErrInvalidPatch
		}
	}

	if rawTitle, ok := courseUpdatePatch["title"]; ok {
		switch title := rawTitle.(type) {
		case string:
			err := c.courseRepository.UpdateTitle(ctx, course.ID, title)
			if err != nil {
				c.logger.Error("failed to update title", "course", course, "newTitle", title, "err", err)
				return err
			}
		default:
			c.logger.Error("title is not a string", "rawTitle", rawTitle)
			return ErrInvalidPatch
		}
	}

	if rawDescription, ok := courseUpdatePatch["description"]; ok {
		switch description := rawDescription.(type) {
		case string:
			err := c.courseRepository.UpdateDescription(ctx, course.ID, description)
			if err != nil {
				c.logger.Error("failed to update description", "course", course, "newDescription", description, "err", err)
				return err
			}
		default:
			c.logger.Error("description is not a string", "rawDescription", rawDescription)
			return ErrInvalidPatch
		}
	}

	if rawSlug, ok := courseUpdatePatch["slug"]; ok {
		switch slug := rawSlug.(type) {
		case string:
			err := c.courseRepository.UpdateSlug(ctx, course.ID, slug)
			if err != nil {
				c.logger.Error("failed to update slug", "course", course, "newSlug", slug, "err", err)
				return err
			}
		default:
			c.logger.Error("slug is not a string", "rawSlug", rawSlug)
			return ErrInvalidPatch
		}
	}

	if rawDifficulty, ok := courseUpdatePatch["difficulty"]; ok {
		switch difficulty := rawDifficulty.(type) {
		case string:
			err := c.courseRepository.UpdateDifficulty(ctx, course.ID, difficulty)
			if err != nil {
				c.logger.Error("failed to update difficulty", "course", course, "newDifficulty", difficulty, "err", err)
				return err
			}
		default:
			c.logger.Error("difficulty is not a string", "rawDifficulty", rawDifficulty)
			return ErrInvalidPatch
		}
	}

	return nil
}

func (c *courseService) UpdateLesson(ctx context.Context, lesson *models.Lesson, lessonUpdatePach map[string]any) error {
	if rawVideoUrl, ok := lessonUpdatePach["videoUrl"]; ok {
		switch videoUrl := rawVideoUrl.(type) {
		case string:
			err := c.lessonRepository.UpdateVideoURL(ctx, lesson.ID, videoUrl)
			if err != nil {
				c.logger.Error("failed to update lesson videoUrl", "lesson", lesson, "newVideoUrl", videoUrl, "err", err)
				return err
			}
		default:
			c.logger.Error("videoUrl is not a string", "rawVideoUrl", rawVideoUrl)
			return ErrInvalidPatch
		}
	}

	if rawDurationSeconds, ok := lessonUpdatePach["durationSeconds"]; ok {
		switch durationSeconds := rawDurationSeconds.(type) {
		case int:
			err := c.lessonRepository.UpdateDurationSeconds(ctx, lesson.ID, durationSeconds)
			if err != nil {
				c.logger.Error("failed to update lesson durationSeconds", "lesson", lesson, "newDurationSeconds", durationSeconds, "err", err)
				return err
			}
		default:
			c.logger.Error("durationSeconds is not a string", "rawDurationSeconds", rawDurationSeconds)
			return ErrInvalidPatch
		}
	}

	if rawTitle, ok := lessonUpdatePach["title"]; ok {
		switch title := rawTitle.(type) {
		case string:
			err := c.lessonRepository.UpdateTitle(ctx, lesson.ID, title)
			if err != nil {
				c.logger.Error("failed to update lesson title", "lesson", lesson, "newtitle", title, "err", err)
				return err
			}
		default:
			c.logger.Error("title is not a string", "rawtitle", rawTitle)
			return ErrInvalidPatch
		}
	}

	if rawDescription, ok := lessonUpdatePach["description"]; ok {
		switch description := rawDescription.(type) {
		case string:
			err := c.lessonRepository.UpdateDescription(ctx, lesson.ID, description)
			if err != nil {
				c.logger.Error("failed to update lesson description", "lesson", lesson, "newdescription", description, "err", err)
				return err
			}
		default:
			c.logger.Error("description is not a string", "rawDescription", rawDescription)
			return ErrInvalidPatch
		}
	}

	if rawLessonNumber, ok := lessonUpdatePach["lessonNumber"]; ok {
		switch lessonNumber := rawLessonNumber.(type) {
		case int:
			err := c.lessonRepository.UpdateLessonNumber(ctx, lesson.ID, lessonNumber)
			if err != nil {
				c.logger.Error("failed to update lesson lessonNumber", "lesson", lesson, "newlessonNumber", lessonNumber, "err", err)
				return err
			}
		default:
			c.logger.Error("lessonNumber is not a string", "rawlessonNumber", rawLessonNumber)
			return ErrInvalidPatch
		}
	}

	return nil
}

func NewCourseService(
	courseRepository repository.CourseRepository,
	lessonRepository repository.LessonRepository,
) CourseService {
	return &courseService{
		courseRepository: courseRepository,
		lessonRepository: lessonRepository,

		logger: slog.Default().With("package", "courseService"),
	}
}
