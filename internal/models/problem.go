package models

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ID uuid.UUID `db:"id" json:"-"`

	Title       string            `db:"title" json:"title"`
	Slug        string            `db:"slug" json:"slug"`
	Description string            `db:"description" json:"description"`
	Difficulty  ProblemDifficulty `json:"difficulty"`

	IsPremium bool `db:"is_premium" json:"isPremium"`
	IsPublic  bool `db:"is_public" json:"isPublic"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type PublicProblem struct {
	*Problem
	Hints        []string                  `json:"hints"`
	CodeSnippets []*ProblemCodeCodeSnippet `json:"codeSnippets"`
	Testcases    []*ProblemTestcase        `json:"testcases"`
}

type ProblemStats struct {
	EasyCount   int `json:"easyCount"`
	MediumCount int `json:"mediumCount"`
	HardCount   int `json:"hardCount"`
}

func NewProblem(
	title string,
	slug string,
	description string,
	difficulty string,

	isPremium bool,
	isPublic bool,
) (*Problem, error) {
	problem := &Problem{
		Title:       title,
		Slug:        slug,
		Description: description,

		IsPremium: isPremium,
		IsPublic:  isPublic,
	}

	var err error
	problem.Difficulty, err = NewProblemDifficulty(difficulty)
	if err != nil {
		return nil, ErrInvalidProblemDifficulty
	}

	return problem, nil
}
