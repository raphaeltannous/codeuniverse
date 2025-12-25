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

type postgresCourseRepository struct {
	db *sql.DB
}

func (p *postgresCourseRepository) Create(ctx context.Context, course *models.Course) (*models.Course, error) {
	query := `
		INSERT INTO courses
			(title, description)
		VALUES
			($1, $2)
		RETURNING id;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		course.Title,
		course.Description,
	)

	err := row.Scan(&course.ID)
	if err != nil {
		return nil, repository.ErrInternalServerError
	}

	course.ThumbnailURL = "default.jpg"
	return course, nil
}

func (p *postgresCourseRepository) Delete(ctx context.Context, courseId uuid.UUID) error {
	query := `
		DELETE FROM courses
		WHERE id = $1;
	`

	result, err := p.db.ExecContext(
		ctx,
		query,
		courseId,
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

func (p *postgresCourseRepository) Get(ctx context.Context, courseId uuid.UUID) (*models.Course, error) {
	query := `
		SELECT id, title, description, thumbnail_url, created_at, updated_at
		FROM courses
		WHERE id = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		courseId,
	)

	course := new(models.Course)
	if err := p.scanCourseFunc(row, course); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrCourseNotFound
		}

		return nil, repository.ErrInternalServerError
	}

	return course, nil
}

func (p *postgresCourseRepository) GetAll(ctx context.Context) ([]*models.Course, error) {
	query := `
		SELECT id, title, description, thumbnail_url, created_at, updated_at
		FROM courses;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query all courses: %w", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := new(models.Course)

		err := p.scanCourseFunc(rows, course)
		if err != nil {
			return nil, fmt.Errorf("failed to scan into course: %w", err)
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (p *postgresCourseRepository) UpdateDescription(ctx context.Context, courseId uuid.UUID, description string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"courses",
		courseId,
		"description",
		description,
	)
}

func (p *postgresCourseRepository) UpdateThumbnailURL(ctx context.Context, courseId uuid.UUID, url string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"courses",
		courseId,
		"thumbnail_url",
		url,
	)
}

func (p *postgresCourseRepository) UpdateTitle(ctx context.Context, courseId uuid.UUID, title string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"courses",
		courseId,
		"title",
		title,
	)
}

func NewPostgresCourseRepository(db *sql.DB) repository.CourseRepository {
	return &postgresCourseRepository{db: db}
}

func (p *postgresCourseRepository) scanCourseFunc(scanner postgresScanner, course *models.Course) error {
	return scanner.Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.ThumbnailURL,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
}
