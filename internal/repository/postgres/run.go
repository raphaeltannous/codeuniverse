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

type postgresRunRepository struct {
	db *sql.DB
}

func NewRunRepository(db *sql.DB) repository.RunRepository {
	return &postgresRunRepository{db: db}
}

func (p *postgresRunRepository) Create(ctx context.Context, run *models.Run) (*models.Run, error) {
	query := `
	INSERT INTO
	  runs (user_id, problem_id, language, code, status)
	VALUES
	  ($1, $2, $3, $4, $5)
	RETURNING
	  id;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		run.UserId,
		run.ProblemId,
		run.Language,
		run.Code,
		run.Status,
	)

	err := row.Scan(&run.ID)
	if err != nil {
		return nil, err
	}

	return run, nil
}

func (p *postgresRunRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return p.updateColumnValue(
		ctx,
		id,
		"status",
		status,
	)
}

func (p *postgresRunRepository) UpdateExecutionTime(ctx context.Context, id uuid.UUID, executionTime float64) error {
	return p.updateColumnValue(
		ctx,
		id,
		"execution_time",
		executionTime,
	)
}

func (p *postgresRunRepository) UpdateMemoryUsage(ctx context.Context, id uuid.UUID, memoryUsage float64) error {
	return p.updateColumnValue(
		ctx,
		id,
		"memory_usage",
		memoryUsage,
	)
}

func (p *postgresRunRepository) UpdateAcceptanceStatus(ctx context.Context, id uuid.UUID, status bool) error {
	return p.updateColumnValue(
		ctx,
		id,
		"is_accepted",
		status,
	)
}

func (p *postgresRunRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Run, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, execution_time, memory_usage, is_accepted, created_at, updated_at
		FROM runs
		WHERE id = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		id,
	)

	run := new(models.Run)
	if err := p.scanRunFunc(row, run); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRunNotFound
		}

		fmt.Println(err)
		return nil, repository.ErrInternalServerError
	}

	return run, nil
}

func (p *postgresRunRepository) updateColumnValue(ctx context.Context, id uuid.UUID, column string, value any) error {
	return updateColumnValue(
		ctx,
		p.db,
		"runs",
		id,
		column,
		value,
	)
}

func (p *postgresRunRepository) scanRunFunc(scanner postgresScanner, run *models.Run) error {
	return scanner.Scan(
		&run.ID,
		&run.UserId,
		&run.ProblemId,
		&run.Language,
		&run.Code,
		&run.Status,
		&run.ExecutionTime,
		&run.MemoryUsage,
		&run.IsAccepted,
		&run.CreatedAt,
		&run.UpdatedAt,
	)
}
