package models

import (
	"time"

	"github.com/google/uuid"
)

type ProblemCategory struct {
	ProblemId  uuid.UUID `db:"problem_id" json:"-"`
	CategoryId uuid.UUID `db:"category_id" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
