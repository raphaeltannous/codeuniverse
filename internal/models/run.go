package models

import (
	"time"

	"github.com/google/uuid"
)

type Run struct {
	ID uuid.UUID `db:"id" json:"id"`

	UserId    uuid.UUID `db:"user_id" json:"-"`
	ProblemId uuid.UUID `db:"problem_id" json:"-"`

	Language string       `db:"language" json:"language"`
	Code     string       `db:"code" json:"code"`
	Status   ResultStatus `db:"status" json:"status"`

	ExecutionTime float64 `db:"execution_time" json:"executionTime"`
	MemoryUsage   float64 `db:"memory_usage" json:"memoryUsage"`

	FailedTestcases []*FailedTestcase `db:"failed_testcases" json:"failedTestcases"`

	StdOut string `db:"stdout" json:"stdOut"`
	StdErr string `db:"sterr" json:"stdErr"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func NewRun(
	userId uuid.UUID,
	problemId uuid.UUID,

	language string,
	code string,
	status ResultStatus,
) *Run {
	return &Run{
		UserId:    userId,
		ProblemId: problemId,

		Language: language,
		Code:     code,
		Status:   status,
	}
}

type RunResult struct {
	Status ResultStatus

	FailedTestcases []*FailedTestcase

	MemoryUsage   float64
	ExecutionTime float64

	StdOut string
	StdErr string
}

func NewRunResult(results *Results) *RunResult {
	runResults := &RunResult{
		FailedTestcases: results.FailedTestcases,
		ExecutionTime:   results.ExecutionTime,
	}

	if results.IsPassed {
		runResults.Status = StatusAccepted
	} else {
		runResults.Status = StatusFailed
	}

	return runResults
}
