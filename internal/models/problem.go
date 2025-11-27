package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ID uuid.UUID `db:"id" json:"-"`

	Title       string `db:"title" json:"title"`
	Slug        string `db:"slug" json:"slug"`
	Description string `db:"description" json:"description"`
	Difficulty  string `db:"difficulty" json:"difficulty"`

	Hints []string `db:"hints" json:"hints"`

	CodeSnippets []byte `db:"code_snippets" json:"codeSnippets"`
	TestCases    []byte `db:"test_cases" json:"testcases"`

	IsPaid   *bool `db:"is_paid" json:"isPaid"`
	IsPublic *bool `db:"is_public" json:"isPublic"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type CodeSnippet struct {
	Code         string `json:"code"`
	Language     string `json:"languageName"`
	LanguageSlug string `json:"languageSlug"`
}

func NewProblem(
	title string,
	description string,
	difficulty string,
	isPaid bool,
	isPublic bool,

	hints []string,

	codeSnippets []CodeSnippet,
	TestCases []string,
) (*Problem, error) {
	codeSnippetsBytes, err := json.Marshal(codeSnippets)
	if err != nil {
		return nil, err
	}

	testCasesBytes, err := json.Marshal(TestCases)
	if err != nil {
		return nil, err
	}

	problem := &Problem{
		Title:       title,
		Slug:        generateSlug(title),
		Description: description,
		Difficulty:  difficulty,
		IsPaid:      &isPaid,
		IsPublic:    &isPublic,

		Hints: hints,

		CodeSnippets: codeSnippetsBytes,
		TestCases:    testCasesBytes,
	}

	return problem, nil
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}
