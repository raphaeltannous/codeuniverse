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

func (p *postgresSubmissionRepository) GetSubmissionsCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM submissions;
	`

	row := p.db.QueryRowContext(ctx, query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get submissions count: %w", err)
	}

	return count, nil
}

func (p *postgresSubmissionRepository) GetSubmissionsLastNDaysCount(ctx context.Context, since int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM submissions
		WHERE created_at >= NOW() - $1::INTERVAL;
	`

	row := p.db.QueryRowContext(ctx, query, fmt.Sprintf("%d days", since))

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get submissions count since %d days: %w", since, err)
	}

	return count, nil
}

func (p *postgresSubmissionRepository) GetPendingSubmissionsCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM submissions
		WHERE status = 'PENDING';
	`

	row := p.db.QueryRowContext(ctx, query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get pending submissions count: %w", err)
	}

	return count, nil
}

func (p *postgresSubmissionRepository) GetAcceptedSubmissionsCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM submissions
		WHERE status = 'ACCEPTED';
	`

	row := p.db.QueryRowContext(ctx, query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get accepted submissions count: %w", err)
	}

	return count, nil
}

func (p *postgresSubmissionRepository) GetRecentSubmissions(ctx context.Context, limit int) ([]*models.SubmissionActivity, error) {
	query := `
        SELECT
            s.id,
            u.username,
            p.title as problem_title,
            p.id as problem_id,
            s.status,
            s.created_at
        FROM submissions s
        JOIN users u ON u.id = s.user_id
        JOIN problems p ON p.id = s.problem_id
        WHERE s.created_at >= NOW() - INTERVAL '24 hours'
        ORDER BY s.created_at DESC
        LIMIT $1
    `

	rows, err := p.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*models.SubmissionActivity
	for rows.Next() {
		submission := new(models.SubmissionActivity)
		if err := rows.Scan(
			&submission.ID,
			&submission.Username,
			&submission.ProblemTitle,
			&submission.ProblemId,
			&submission.Status,
			&submission.CreatedAt,
		); err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (p *postgresSubmissionRepository) GetDailySubmissions(ctx context.Context, since int) ([]*models.DailySubmissions, error) {
	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as submissions,
			COUNT(CASE WHEN status = 'ACCEPTED' THEN 1 END) as accepted
		FROM submissions
		WHERE created_at >= NOW() - $1::interval
		GROUP BY DATE(created_at)
		ORDER BY date;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		fmt.Sprintf("%d days", since),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily submissions: %w", err)
	}
	defer rows.Close()

	var dailySubmissions []*models.DailySubmissions
	for rows.Next() {
		dailySubmission := new(models.DailySubmissions)
		if err := rows.Scan(
			&dailySubmission.Date,
			&dailySubmission.Submissions,
			&dailySubmission.Accepted,
		); err != nil {
			return nil, fmt.Errorf("failed to scan into dailySubmission: %w", err)
		}

		dailySubmissions = append(dailySubmissions, dailySubmission)
	}

	return dailySubmissions, nil
}

func (p *postgresSubmissionRepository) GetDailySubmissionsHours(ctx context.Context, since int) ([]*models.DailySubmissions, error) {
	query := `
		SELECT
			DATE_TRUNC('hour', created_at) as date,
			COUNT(*) as submissions,
			COUNT(CASE WHEN status = 'ACCEPTED' THEN 1 END) as accepted
		FROM submissions
		WHERE created_at >= NOW() - $1::interval
		GROUP BY DATE_TRUNC('hour', created_at)
		ORDER BY date;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		fmt.Sprintf("%d hours", since),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily submissions: %w", err)
	}
	defer rows.Close()

	var dailySubmissions []*models.DailySubmissions
	for rows.Next() {
		dailySubmission := new(models.DailySubmissions)
		if err := rows.Scan(
			&dailySubmission.Date,
			&dailySubmission.Submissions,
			&dailySubmission.Accepted,
		); err != nil {
			return nil, fmt.Errorf("failed to scan into dailySubmission: %w", err)
		}

		dailySubmissions = append(dailySubmissions, dailySubmission)
	}

	return dailySubmissions, nil
}

func (p *postgresSubmissionRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, execution_time, memory_usage, is_accepted, created_at, updated_at
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
		SELECT id, user_id, problem_id, language, code, status, execution_time, memory_usage, is_accepted, created_at, updated_at
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

func (p *postgresSubmissionRepository) GetSubmissionsStats(ctx context.Context, userId uuid.UUID) (*models.SubmissionStats, error) {
	query := `
		SELECT
		    COUNT(*) AS total_submissions,
		    COUNT(*) FILTER (WHERE is_accepted = true) AS accepted_submissions,

		    COUNT(DISTINCT problem_id) FILTER (WHERE is_accepted = true) AS problems_solved,

		    COUNT(DISTINCT problem_id) FILTER (
		        WHERE is_accepted = true AND p.difficulty = 'Easy'
		    ) AS easy_solved,

		    COUNT(DISTINCT problem_id) FILTER (
		        WHERE is_accepted = true AND p.difficulty = 'Medium'
		    ) AS medium_solved,

		    COUNT(DISTINCT problem_id) FILTER (
		        WHERE is_accepted = true AND p.difficulty = 'Hard'
		    ) AS hard_solved
		FROM submissions s
		JOIN problems p ON p.id = s.problem_id
		WHERE s.user_id = $1;
	`

	row := p.db.QueryRowContext(
		ctx,
		query,
		userId,
	)

	submissionStats := new(models.SubmissionStats)
	err := row.Scan(
		&submissionStats.TotalSubmissions,
		&submissionStats.AcceptedSubmissions,
		&submissionStats.ProblemsSolved,
		&submissionStats.EasySolved,
		&submissionStats.MediumSolved,
		&submissionStats.HardSolved,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}

		return nil, repository.ErrInternalServerError
	}

	return submissionStats, nil
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
		&submission.UserId,
		&submission.ProblemId,
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
