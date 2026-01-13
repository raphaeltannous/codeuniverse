package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (p *postgresSubmissionRepository) UpdateFailedTestcases(ctx context.Context, id uuid.UUID, failedTestcases []*models.FailedTestcase) error {
	if failedTestcases == nil {
		return nil
	}

	return p.updateColumnValue(
		ctx,
		id,
		"failed_testcases",
		failedTestcases,
	)
}

func (p *postgresSubmissionRepository) UpdateStderr(ctx context.Context, id uuid.UUID, stderr string) error {
	return p.updateColumnValue(
		ctx,
		id,
		"stderr",
		stderr,
	)
}

func (p *postgresSubmissionRepository) UpdateStdout(ctx context.Context, id uuid.UUID, stdout string) error {
	return p.updateColumnValue(
		ctx,
		id,
		"stdout",
		stdout,
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
		WHERE status = 'Accepted';
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
			COUNT(CASE WHEN status = 'Accepted' THEN 1 END) as accepted
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
			COUNT(CASE WHEN status = 'Accepted' THEN 1 END) as accepted
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
		SELECT
			id,

			user_id,
			problem_id,

			language,
			code,
			status,

			execution_time,
			memory_usage,

			failed_testcases,
			stdout,
			stderr,

			created_at,
			updated_at
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
		SELECT
			id,

			user_id,
			problem_id,

			language,
			code,
			status,

			execution_time,
			memory_usage,

			failed_testcases,
			stdout,
			stderr,

			created_at,
			updated_at
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

func (p *postgresSubmissionRepository) GetSolvedProblems(ctx context.Context, userId uuid.UUID) ([]string, error) {
	query := `
		SELECT
			p.slug
		FROM submissions s
		JOIN problems p ON p.id = s.problem_id
		WHERE s.status = 'Accepted' AND s.user_id = $1;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		userId,
	)
	if err != nil {
		return []string(nil), fmt.Errorf("failed to get rows for solved problems: %w", err)
	}

	var solvedProblems []string
	for rows.Next() {
		var slug string

		if err := rows.Scan(&slug); err != nil {
			return []string(nil), fmt.Errorf("failed to scan slug from solved prolbems: %w", err)
		}

		solvedProblems = append(solvedProblems, slug)
	}

	return solvedProblems, nil
}
func (p *postgresSubmissionRepository) GetSubmissionsStats(ctx context.Context, userId uuid.UUID) (*models.SubmissionStats, error) {
	query := `
		SELECT
		    COUNT(*) AS total_submissions,
		    COUNT(*) FILTER (WHERE status = 'Accepted') AS accepted_submissions,

		    COUNT(DISTINCT problem_id) FILTER (WHERE status = 'Accepted') AS problems_solved,

		    COUNT(DISTINCT problem_id) FILTER (
		        WHERE status = 'Accepted' AND p.difficulty = 'Easy'
		    ) AS easy_solved,

		    COUNT(DISTINCT problem_id) FILTER (
		        WHERE status = 'Accepted' AND p.difficulty = 'Medium'
		    ) AS medium_solved,

		    COUNT(DISTINCT problem_id) FILTER (
		        WHERE status = 'Accepted' AND p.difficulty = 'Hard'
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
	var failedTestcasesJson []byte
	var status string

	err := scanner.Scan(
		&submission.ID,

		&submission.UserId,
		&submission.ProblemId,

		&submission.Language,
		&submission.Code,
		&status,

		&submission.ExecutionTime,
		&submission.MemoryUsage,

		&failedTestcasesJson,
		&submission.StdOut,
		&submission.StdErr,

		&submission.CreatedAt,
		&submission.UpdatedAt,
	)
	if err != nil {
		return err
	}

	var failedTestcases []*models.FailedTestcase

	err = json.Unmarshal(failedTestcasesJson, &failedTestcases)
	if err != nil {
		return err
	}

	submission.FailedTestcases = failedTestcases

	submission.Status, err = models.ParseResultStatus(status)
	if err != nil {
		return err
	}

	return nil
}

func NewSubmissionRepository(db *sql.DB) repository.SubmissionRepository {
	return &postgresSubmissionRepository{db: db}
}
