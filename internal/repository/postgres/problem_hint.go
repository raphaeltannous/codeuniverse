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

type postgresProblemHintRepository struct {
	db *sql.DB
}

func (p *postgresProblemHintRepository) Create(
	ctx context.Context,
	hint *models.ProblemHint,
) error {
	query := `
		INSERT INTO problem_hints (problem_id, hint)
		VALUES ($1, $2);
	`

	_, err := p.db.ExecContext(
		ctx,
		query,
		hint.ProblemId,
		hint.Hint,
	)
	if err != nil {
		return fmt.Errorf("failed to insert into problem_hints: %w", err)
	}

	return nil
}

func (p *postgresProblemHintRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
		DELETE FROM problem_hints
		WHERE id = $1;
	`

	_, err := p.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete problem_hint: %w", err)
	}

	return nil
}

func (p *postgresProblemHintRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*models.ProblemHint, error) {
	query := `
		SELECT
			id,
			problem_id,
			hint,
			created_at,
			updated_at
		FROM problem_hints
		WHERE id = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		id,
	)

	hint := new(models.ProblemHint)
	if err := p.scanProblemHintFunc(row, hint); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProblemHintNotFound
		}

		return nil, fmt.Errorf("failed to scan into hint: %w", err)
	}

	return hint, nil
}

func (p *postgresProblemHintRepository) GetHints(
	ctx context.Context,
	problemId uuid.UUID,
) ([]*models.ProblemHint, error) {
	query := `
		SELECT
			id,
			problem_id,
			hint,
			created_at,
			updated_at
		FROM problem_hints
		WHERE problem_id = $1;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		problemId,
	)
	if err != nil {
		return []*models.ProblemHint(nil), fmt.Errorf("failed to get hint rows: %w", err)
	}

	var hints []*models.ProblemHint
	for rows.Next() {
		hint := new(models.ProblemHint)

		if err := p.scanProblemHintFunc(rows, hint); err != nil {
			return []*models.ProblemHint(nil), fmt.Errorf("failed to scan into hint: %w", err)
		}

		hints = append(hints, hint)
	}

	if hints == nil {
		hints = []*models.ProblemHint{}
	}

	return hints, nil
}

func (p *postgresProblemHintRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	hint string,
) error {
	query := `
		UPDATE problem_hints
		SET hint = $1
		WHERE id = $2;
	`

	_, err := p.db.ExecContext(
		ctx,
		query,
		hint,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update hint: %w", err)
	}

	return nil
}

func (p *postgresProblemHintRepository) scanProblemHintFunc(
	scanner postgresScanner,
	hint *models.ProblemHint,
) error {
	return scanner.Scan(
		&hint.ID,
		&hint.ProblemId,
		&hint.Hint,
		&hint.CreatedAt,
		&hint.UpdatedAt,
	)

}

func NewProblemHintRepository(
	db *sql.DB,
) repository.ProblemHintRepository {
	return &postgresProblemHintRepository{
		db: db,
	}
}
