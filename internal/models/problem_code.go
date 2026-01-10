package models

import "encoding/json"

type ProblemCode struct {
	CodeSnippet  string          `json:"codeSnippet"`
	Checker      string          `json:"checker"`
	Driver       string          `json:"driver"`
	RuntimeFiles []string        `json:"-"`
	IsPublic     bool            `json:"isPublic"`
	Language     ProblemLanguage `json:"language"`
}

func (p ProblemCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"codeSnippet":  p.CodeSnippet,
		"checker":      p.Checker,
		"driver":       p.Driver,
		"isPublic":     p.IsPublic,
		"languageName": p.Language.String(),
		"languageSlug": p.Language.Slug(),
	})
}

type ProblemCodeCodeSnippet struct {
	Code     string          `json:"code"`
	Language ProblemLanguage `json:"language"`
}

func (p ProblemCodeCodeSnippet) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"code":         p.Code,
		"languageName": p.Language.String(),
		"languageSlug": p.Language.Slug(),
	})
}

type ProblemCodeMetadata struct {
	IsPublic bool `json:"isPublic"`
}

func DefaultProblemCodeMetadata() *ProblemCodeMetadata {
	return &ProblemCodeMetadata{
		IsPublic: false,
	}
}
