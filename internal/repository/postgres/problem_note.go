package postgres

import (
	"context"
	"database/sql"
	"errors"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgresProblemNoteRepository struct {
	db *sql.DB
}

func (p *postgresProblemNoteRepository) Create(ctx context.Context, note *models.ProblemNote) (*models.ProblemNote, error) {
	query := `
		INSERT INTO problem_notes (user_id, problem_id, markdown)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		note.UserId,
		note.ProblemId,
		note.Markdown,
	)

	err := row.Scan(&note.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return note, repository.ErrProblemNoteAlreadyExists
		}

		return note, repository.ErrInternalServerError
	}

	return note, nil
}

func (p *postgresProblemNoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM problem_notes
		where id = $1;
	`

	result, err := p.db.ExecContext(
		ctx,
		query,
		id,
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

func (p *postgresProblemNoteRepository) Get(ctx context.Context, userId uuid.UUID, problemId uuid.UUID) (*models.ProblemNote, error) {
	query := `
		SELECT id, user_id, problem_id, markdown, created_at, updated_at
		FROM problem_notes
		WHERE user_id = $1 AND problem_id = $2;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		userId,
		problemId,
	)

	note := new(models.ProblemNote)
	if err := p.scanProblemNoteFunc(row, note); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProblemNoteNotFound
		}

		return nil, repository.ErrInternalServerError
	}

	return note, nil
}

func (p *postgresProblemNoteRepository) UpdateMarkdown(ctx context.Context, note *models.ProblemNote, markdown string) error {
	return updateColumnValue(
		ctx,
		p.db,
		"problem_notes",
		note.ID,
		"markdown",
		markdown,
	)
}

func (p *postgresProblemNoteRepository) scanProblemNoteFunc(scanner postgresScanner, note *models.ProblemNote) error {
	var markdown *string

	err := scanner.Scan(
		&note.ID,
		&note.UserId,
		&note.ProblemId,
		&markdown,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if markdown == nil {
		note.Markdown = ""
	} else {
		note.Markdown = *markdown
	}

	return nil
}

func NewProblemNoteRepository(db *sql.DB) repository.ProblemNoteRepository {
	return &postgresProblemNoteRepository{db: db}
}
