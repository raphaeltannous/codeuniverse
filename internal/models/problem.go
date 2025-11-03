package models

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title       string `db:"title" json:"title"`
	Slug        string `db:"slug" json:"slug"`
	Number      int    `db:"number" json:"number"`
	Description string `db:"description" json:"description"`
	Difficulty  string `db:"difficulty" json:"difficulty"`

	Tags  []string `db:"tags" json:"tags"`
	Hints []string `db:"hints" json:"hints"`

	CodeSnippets []byte `db:"code_snippets" json:"codeSnippets"`
	TestCases    []byte `db:"test_cases" json:"testCases"`

	Likes    *int `db:"likes" json:"likes"`
	Dislikes *int `db:"dislikes" json:"dislikes"`

	IsPaid   *bool `db:"is_paid" json:"isPaid"`
	IsPublic *bool `db:"is_public" json:"isPublic"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
