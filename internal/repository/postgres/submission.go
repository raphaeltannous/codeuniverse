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

type postgresSubmissionRepository struct {
	db *sql.DB
}

func (p *postgresSubmissionRepository) Create(ctx context.Context, submission *models.Submission) (*models.Submission, error) {
	query := `
		INSERT INTO
			submissions (user_id, problem_id, language, code, status)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		submission.UserId,
		submission.ProblemId,
		submission.Language,
		submission.Code,
		submission.Status,
	)

	err := row.Scan(&submission.ID)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (p *postgresSubmissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

func (p *postgresSubmissionRepository) UpdateAcceptanceStatus(ctx context.Context, id uuid.UUID, status bool) error {
	return p.updateColumnValue(
		ctx,
		id,
		"is_accepted",
		status,
	)
}

func (p *postgresSubmissionRepository) UpdateExecutionTime(ctx context.Context, id uuid.UUID, executionTime float64) error {
	return p.updateColumnValue(
		ctx,
		id,
		"execution_time",
		executionTime,
	)
}

func (p *postgresSubmissionRepository) UpdateMemoryUsage(ctx context.Context, id uuid.UUID, memoryUsage float64) error {
	return p.updateColumnValue(
		ctx,
		id,
		"memory_usage",
		memoryUsage,
	)
}

func (p *postgresSubmissionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return p.updateColumnValue(
		ctx,
		id,
		"status",
		status,
	)
}

func (p *postgresSubmissionRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	query := `
		SELECT id, language, code, status, execution_time, memory_usage, is_accepted, created_at, updated_at
		FROM submissions
		WHERE id = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		id,
	)

	submission := new(models.Submission)
	if err := p.scanSubmissionFunc(row, submission); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrSubmissionNotFound
		}

		return nil, repository.ErrInternalServerError
	}

	return submission, nil
}

func (p *postgresSubmissionRepository) GetProblemSubmissions(ctx context.Context, userId uuid.UUID, problemId uuid.UUID) ([]*models.Submission, error) {
	query := `
		SELECT id, language, code, status, execution_time, memory_usage, is_accepted, created_at, updated_at
		FROM submissions
		WHERE user_id = $1 AND problem_id = $2;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		userId,
		problemId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query all problems: %w", err)
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		submission := new(models.Submission)

		err := p.scanSubmissionFunc(rows, submission)
		if err != nil {
			return nil, fmt.Errorf("faild to scan into problem: %w", err)
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (p *postgresSubmissionRepository) updateColumnValue(ctx context.Context, id uuid.UUID, column string, value any) error {
	return updateColumnValue(
		ctx,
		p.db,
		"submissions",
		id,
		column,
		value,
	)
}

func (p *postgresSubmissionRepository) scanSubmissionFunc(scanner postgresScanner, submission *models.Submission) error {
	return scanner.Scan(
		&submission.ID,
		&submission.Language,
		&submission.Code,
		&submission.Status,
		&submission.ExecutionTime,
		&submission.MemoryUsage,
		&submission.IsAccepted,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	)
}

func NewSubmissionRepository(db *sql.DB) repository.SubmissionRepository {
	return &postgresSubmissionRepository{db: db}
}
