package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type postgresCourseProgressRepository struct {
	db *sql.DB
}

func (p *postgresCourseProgressRepository) Get(ctx context.Context, courseId uuid.UUID, userId uuid.UUID) ([]*models.CourseLessonProgress, error) {
	query := `
		SELECT user_id, course_id, lesson_id, is_completed, created_at, updated_at
		FROM course_progress
		WHERE course_id = $1 AND user_id = $2;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		courseId,
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrCourseProgressNotFound
		}

		return nil, fmt.Errorf("failed to get course progress for user: %w", err)
	}

	var courseProgress []*models.CourseLessonProgress
	for rows.Next() {
		progress := new(models.CourseLessonProgress)

		if err := rows.Scan(
			&progress.UserId,
			&progress.CourseId,
			&progress.LessonId,
			&progress.IsCompleted,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan into progress: %w", err)
		}

		courseProgress = append(courseProgress, progress)
	}

	return courseProgress, nil
}

func (p *postgresCourseProgressRepository) Save(ctx context.Context, courseProgress *models.CourseLessonProgress) error {
	query := `
		INSERT INTO course_progress (user_id, course_id, lesson_id, is_completed)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, lesson_id)
		DO UPDATE SET
			is_completed = EXCLUDED.is_completed;
	`

	_, err := p.db.ExecContext(
		ctx,
		query,
		courseProgress.UserId,
		courseProgress.CourseId,
		courseProgress.LessonId,
		courseProgress.IsCompleted,
	)

	return err
}

func NewCourseProgressRepository(db *sql.DB) repository.CourseProgressRepository {
	return &postgresCourseProgressRepository{db: db}
}
