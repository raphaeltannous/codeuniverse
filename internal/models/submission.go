package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
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

func NewSubmission(
	userId uuid.UUID,
	problemId uuid.UUID,

	language string,
	code string,
	status ResultStatus,
) *Submission {
	return &Submission{
		UserId:    userId,
		ProblemId: problemId,

		Language: language,
		Code:     code,
		Status:   status,
	}
}
