package filesystem

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
)

var (
	ErrBaseDirectoryIsNotAbsolute = errors.New("baseDirectory path is not absolute path")
)

type filesystemProblemCodeRepository struct {
	baseDirectory string
}

func (f *filesystemProblemCodeRepository) SaveCodeSnippet(ctx context.Context, problem *models.Problem, snippet *models.CodeSnippet) error {
	filename := snippet.Language.CodeSnippetFilename()

	return f.saveLanguageFile(problem, &snippet.Language, filename, []byte(snippet.Code))
}

func (f *filesystemProblemCodeRepository) saveLanguageFile(
	problem *models.Problem,
	language *models.ProblemLanguage,
	filename string,
	content []byte,
) error {
	languagePath, err := f.getLanguagePath(problem, language)
	if err != nil {
		return err
	}

	filePath := filepath.Join(languagePath, filename)
	return os.WriteFile(filePath, content, 0640)
}

func (f *filesystemProblemCodeRepository) getProblemPath(problem *models.Problem) (string, error) {
	problemPath := filepath.Join(f.baseDirectory, problem.ID.String())

	if _, err := os.Stat(problemPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(problemPath, 0750); err != nil {
			return "", err
		}
	}

	return problemPath, nil
}

func (f *filesystemProblemCodeRepository) getLanguagePath(problem *models.Problem, language *models.ProblemLanguage) (string, error) {
	problemPath, err := f.getProblemPath(problem)
	if err != nil {
		return "", err
	}

	languagePath := filepath.Join(problemPath, language.Slug())
	if _, err := os.Stat(languagePath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(languagePath, 0750); err != nil {
			return "", err
		}
	}

	return languagePath, nil
}

func NewFilesystemProblemCodeRepository(baseDirectory string) (repository.ProblemCodeRepository, error) {
	if !filepath.IsAbs(baseDirectory) {
		return nil, ErrBaseDirectoryIsNotAbsolute
	}

	return &filesystemProblemCodeRepository{
		baseDirectory: baseDirectory,
	}, nil
}
