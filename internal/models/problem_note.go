package models

import (
	"time"

	"github.com/google/uuid"
)

type ProblemNote struct {
	ID uuid.UUID `db:"id" json:"-"`

	UserId    uuid.UUID `db:"user_id" json:"-"`
	ProblemId uuid.UUID `db:"problem_id" json:"-"`

	Markdown string `db:"markdown" json:"markdown"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func NewProblemNote(
	userId uuid.UUID,
	problemId uuid.UUID,
	markdown string,
) *ProblemNote {
	return &ProblemNote{
		UserId:    userId,
		ProblemId: problemId,
		Markdown:  markdown,
	}
}
