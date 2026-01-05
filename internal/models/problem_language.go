package models

import (
	"errors"
	"fmt"
	"strings"
)

type ProblemLanguage int

const (
	LanguageGo ProblemLanguage = iota + 1
	LanguagePython
	LanguageCpp
	LanguageTypescript
	LanguageJavascript
)

var (
	ErrInvalidProblemLanguage = errors.New("invalid problem language")
)

func NewProblemLanguage(language string) (ProblemLanguage, error) {
	switch strings.ToLower(language) {
	case "go":
		return LanguageGo, nil
	case "python":
		return LanguagePython, nil
	case "cpp":
		return LanguageCpp, nil
	case "typescript":
		return LanguageTypescript, nil
	case "javascript":
		return LanguageJavascript, nil
	default:
		return 0, ErrInvalidProblemLanguage
	}
}

func (p ProblemLanguage) String() string {
	switch p {
	case LanguageGo:
		return "Go"
	case LanguagePython:
		return "Python"
	case LanguageCpp:
		return "Cpp"
	case LanguageTypescript:
		return "Typescript"
	case LanguageJavascript:
		return "Javascript"
	default:
		return "Unknown"
	}
}

func (p ProblemLanguage) Slug() string {
	return fmt.Sprintf("%s", strings.ToLower(p.String()))
}

func (p ProblemLanguage) FileExtension() string {
	switch p {
	case LanguageGo:
		return ".go"
	case LanguagePython:
		return ".py"
	case LanguageCpp:
		return ".cpp"
	case LanguageTypescript:
		return ".ts"
	case LanguageJavascript:
		return ".js"
	default:
		return ".txt"
	}
}

func (p ProblemLanguage) CodeSnippetFilename() string {
	filename := "unknown"

	switch p {
	case LanguageGo:
		filename = "main"
	case LanguagePython:
		filename = "main"
	case LanguageCpp:
		filename = "main"
	case LanguageTypescript:
		filename = "main"
	case LanguageJavascript:
		filename = "main"
	default:
		filename = "unknown"
	}

	return filename + p.FileExtension()
}
