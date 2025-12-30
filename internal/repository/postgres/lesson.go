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

type postgresLessonRepository struct {
	db *sql.DB
}

func (p *postgresLessonRepository) GetLessonCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM lessons;
	`

	row := p.db.QueryRowContext(ctx, query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *postgresLessonRepository) GetLessonCountForCourse(ctx context.Context, courseId uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM lessons
		WHERE course_id = $1;
	`

	row := p.db.QueryRowContext(ctx, query, courseId)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *postgresLessonRepository) Create(ctx context.Context, courseId uuid.UUID, lesson *models.Lesson) (*models.Lesson, error) {
	query := `
		INSERT INTO lessons
			(course_id, title, description, lesson_number)
		VALUES
			($1, $2, $3, $4)
		RETURNING id;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		courseId,
		lesson.Title,
		lesson.Description,
		lesson.LessonNumber,
	)

	err := row.Scan(&lesson.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create lesson: %w", err)
	}

	return lesson, nil
}

func (p *postgresLessonRepository) Delete(ctx context.Context, lessonId uuid.UUID) error {
	query := `
		DELETE FROM lessons
		WHERE id = $1;
	`

	result, err := p.db.ExecContext(
		ctx,
		query,
		lessonId,
	)

	if err != nil {
		return repository.ErrInternalServerError
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return repository.ErrInternalServerError
	}

	if rows != 1 {
		return repository.ErrInternalServerError
	}

	return nil
}

func (p *postgresLessonRepository) Get(ctx context.Context, lessonId uuid.UUID) (*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, description, video_url, duration_seconds, lesson_number, created_at, updated_at
		FROM lessons
		WHERE id = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		lessonId,
	)

	lesson := new(models.Lesson)
	if err := p.scanLessonFunc(row, lesson); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrLessonNotFound
		}

		return nil, repository.ErrInternalServerError
	}

	return lesson, nil
}

func (p *postgresLessonRepository) GetAllByCourse(ctx context.Context, courseId uuid.UUID) ([]*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, description, video_url, duration_seconds, lesson_number, created_at, updated_at
		FROM lessons
		WHERE course_id = $1;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		courseId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query all lessons by courseId: %w", err)
	}
	defer rows.Close()

	var lessons []*models.Lesson
	for rows.Next() {
		lesson := new(models.Lesson)

		err := p.scanLessonFunc(rows, lesson)
		if err != nil {
			return nil, fmt.Errorf("failed to scan into lesson: %w", err)
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (p *postgresLessonRepository) UpdateDescription(ctx context.Context, lessonId uuid.UUID, description string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"lessons",
		lessonId,
		"description",
		description,
	)
}

func (p *postgresLessonRepository) UpdateDurationSeconds(ctx context.Context, lessonId uuid.UUID, duration int) error {
	return updateColumnValue(
		ctx,
		p.db,
		"lessons",
		lessonId,
		"duration_seconds",
		duration,
	)
}

func (p *postgresLessonRepository) UpdateLessonNumber(ctx context.Context, lessonId uuid.UUID, lessonNumber int) error {
	return updateColumnValue(
		ctx,
		p.db,
		"lessons",
		lessonId,
		"lesson_number",
		lessonNumber,
	)
}

func (p *postgresLessonRepository) UpdateTitle(ctx context.Context, lessonId uuid.UUID, title string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"lessons",
		lessonId,
		"title",
		title,
	)
}

func (p *postgresLessonRepository) UpdateVideoURL(ctx context.Context, lessonId uuid.UUID, url string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"lessons",
		lessonId,
		"video_url",
		url,
	)
}

func NewPostgresLessonRepository(db *sql.DB) repository.LessonRepository {
	return &postgresLessonRepository{db: db}
}

func (p *postgresLessonRepository) scanLessonFunc(scanner postgresScanner, lesson *models.Lesson) error {
	return scanner.Scan(
		&lesson.ID,
		&lesson.CourseId,
		&lesson.Title,
		&lesson.Description,
		&lesson.VideoURL,
		&lesson.DurationSeconds,
		&lesson.LessonNumber,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
}
