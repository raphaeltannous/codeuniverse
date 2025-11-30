package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID uuid.UUID `db:"id" json:"-"`

	UserId    uuid.UUID `db:"user_id" json:"-"`
	ProblemId uuid.UUID `db:"problem_id" json:"-"`

	Language string `db:"language" json:"language"`
	Code     string `db:"code" json:"code"`
	Status   string `db:"status" json:"status"`

	ExecutionTime float64 `db:"execution_time" json:"executionTime"`
	MemoryUsage   float64 `db:"memory_usage" json:"memoryUsage"`

	IsAccepted bool `db:"is_accepted" json:"isAccepted"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
