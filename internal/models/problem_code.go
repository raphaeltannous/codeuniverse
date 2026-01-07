package models

import "encoding/json"

type ProblemCode struct {
	CodeSnippet string          `json:"codeSnippet"`
	Checker     string          `json:"checker"`
	Driver      string          `json:"driver"`
	IsPublic    bool            `json:"isPublic"`
	Language    ProblemLanguage `json:"language"`
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
	Code         string `json:"code"`
	LanguageName string `json:"languageName"`
	LanguageSlug string `json:"languageSlug"`
}

type ProblemCodeMetadata struct {
	IsPublic bool `json:"isPublic"`
}

func DefaultProblemCodeMetadata() *ProblemCodeMetadata {
	return &ProblemCodeMetadata{
		IsPublic: false,
	}
}
