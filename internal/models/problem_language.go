package models

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidProblemLanguage = errors.New("invalid problem language")
)

type ProblemLanguage int

func (p ProblemLanguage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"languageName": p.String(),
		"languageSlug": p.Slug(),
	})
}

const (
	LanguageGo ProblemLanguage = iota + 1
	LanguagePython
	LanguageCpp
	LanguageTypescript
	LanguageJavascript
	LanguageJava
	LanguageRuby

	LanguageEnd
)

type problemLanguageMetadata struct {
	name                string
	slug                string
	fileExtension       string
	codeSnippetFilename string
	driverFilename      string
	checkerFilename     string
}

var problemLanguages = map[ProblemLanguage]*problemLanguageMetadata{
	LanguageGo: {
		name:                "Go",
		slug:                "go",
		fileExtension:       ".go",
		codeSnippetFilename: "main",
		driverFilename:      "driver",
		checkerFilename:     "main_test",
	},
	LanguagePython: {
		name:                "Python",
		slug:                "python",
		fileExtension:       ".py",
		codeSnippetFilename: "main",
		driverFilename:      "driver",
		checkerFilename:     "main_test",
	},
	LanguageCpp: {
		name:                "Cpp",
		slug:                "cpp",
		fileExtension:       ".cpp",
		codeSnippetFilename: "main",
		driverFilename:      "driver",
		checkerFilename:     "main_test",
	},
	LanguageTypescript: {
		name:                "Typescript",
		slug:                "typescript",
		fileExtension:       ".ts",
		codeSnippetFilename: "main",
		driverFilename:      "driver",
		checkerFilename:     "main_test",
	},
	LanguageJavascript: {
		name:                "Javascript",
		slug:                "javascript",
		fileExtension:       ".js",
		codeSnippetFilename: "main",
		driverFilename:      "driver",
		checkerFilename:     "main_test",
	},
	LanguageJava: {
		name:                "Java",
		slug:                "java",
		fileExtension:       ".java",
		codeSnippetFilename: "Main",
		driverFilename:      "Driver",
		checkerFilename:     "MainTest",
	},
	LanguageRuby: {
		name:                "Ruby",
		slug:                "ruby",
		fileExtension:       ".rb",
		codeSnippetFilename: "main",
		driverFilename:      "driver",
		checkerFilename:     "main_test",
	},
}

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
	case "java":
		return LanguageJava, nil
	case "ruby":
		return LanguageRuby, nil
	default:
		return 0, ErrInvalidProblemLanguage
	}
}

func (p ProblemLanguage) String() string {
	if language, ok := problemLanguages[p]; ok {
		return language.name
	}

	return "Unknown"
}

func (p ProblemLanguage) Slug() string {
	if language, ok := problemLanguages[p]; ok {
		return language.slug
	}

	return "unknown"
}

func (p ProblemLanguage) FileExtension() string {
	if language, ok := problemLanguages[p]; ok {
		return language.fileExtension
	}

	return ".txt"
}

func (p ProblemLanguage) CodeSnippetFilename() string {
	if language, ok := problemLanguages[p]; ok {
		return language.codeSnippetFilename + language.fileExtension
	}

	return "unknown"
}

func (p ProblemLanguage) DriverFilename() string {
	if language, ok := problemLanguages[p]; ok {
		return language.driverFilename + language.fileExtension
	}

	return "unknown"
}

func (p ProblemLanguage) CheckerFilename() string {
	if language, ok := problemLanguages[p]; ok {
		return language.checkerFilename + language.fileExtension
	}

	return "unknown"
}
