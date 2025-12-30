package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgresCourseRepository struct {
	db *sql.DB
}

func (p *postgresCourseRepository) Create(ctx context.Context, course *models.Course) (*models.Course, error) {
	query := `
		INSERT INTO courses
			(title, slug, description, difficulty)
		VALUES
			($1, $2, $3, $4)
		RETURNING id;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		course.Title,
		course.Slug,
		course.Description,
		course.Difficulty,
	)

	err := row.Scan(&course.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, repository.ErrCourseAlreadyExists
		}
		return nil, fmt.Errorf("failed to create course: %w", err)
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

func (p *postgresCourseRepository) GetBySlug(ctx context.Context, slug string) (*models.Course, error) {
	query := `
		SELECT id, title, slug, description, difficulty, is_published, thumbnail_url, created_at, updated_at
		FROM courses
		WHERE slug = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		slug,
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
		SELECT id, title, slug, description, difficulty, is_published, thumbnail_url, created_at, updated_at
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

func (p *postgresCourseRepository) UpdateIsPublished(ctx context.Context, courseId uuid.UUID, status bool) error {
	return updateColumnValue(
		ctx,
		p.db,
		"courses",
		courseId,
		"is_published",
		status,
	)
}

func (p *postgresCourseRepository) UpdateDifficulty(ctx context.Context, courseId uuid.UUID, difficulty string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"courses",
		courseId,
		"difficulty",
		difficulty,
	)
}

func (p *postgresCourseRepository) UpdateSlug(ctx context.Context, courseId uuid.UUID, slug string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"courses",
		courseId,
		"slug",
		slug,
	)
}

func NewPostgresCourseRepository(db *sql.DB) repository.CourseRepository {
	return &postgresCourseRepository{db: db}
}

func (p *postgresCourseRepository) scanCourseFunc(scanner postgresScanner, course *models.Course) error {
	return scanner.Scan(
		&course.ID,
		&course.Title,
		&course.Slug,
		&course.Description,
		&course.Difficulty,
		&course.IsPublished,
		&course.ThumbnailURL,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
}
