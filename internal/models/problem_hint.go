package models

import (
	"time"

	"github.com/google/uuid"
)

type ProblemHint struct {
	ID uuid.UUID `db:"id" json:"id"`

	ProblemId uuid.UUID `db:"problem_id" json:"-"`

	Hint string `db:"hint" json:"hint"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
