package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (f *filesystemProblemCodeRepository) DeleteProblem(ctx context.Context, problem *models.Problem) error {
	problemPath, err := f.getProblemPath(problem)
	if err != nil {
		return fmt.Errorf("failed to delete problemPath: %w", err)
	}

	return os.RemoveAll(problemPath)
}

func (f *filesystemProblemCodeRepository) GetProblemCodeConfig(
	ctx context.Context,
	problem *models.Problem,
) (*models.ProblemCodeConfig, error) {
	problemPath, err := f.getProblemPath(problem)
	if err != nil {
		return nil, fmt.Errorf("failed to get problemPath: %w", err)
	}
	configPath := filepath.Join(problemPath, "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			config := models.DefaultProblemCodeConfig()
			if err := f.SaveProblemCodeConfig(ctx, problem, config); err != nil {
				return nil, fmt.Errorf("failed to save default config: %w", err)
			}

			return config, nil
		}

		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config *models.ProblemCodeConfig
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}

func (f *filesystemProblemCodeRepository) SaveProblemCodeConfig(
	ctx context.Context,
	problem *models.Problem,
	config *models.ProblemCodeConfig,
) error {
	jsonData, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal indent config: %w", err)
	}

	problemPath, err := f.getProblemPath(problem)
	if err != nil {
		return fmt.Errorf("failed to get problemPath: %w", err)
	}
	configPath := filepath.Join(problemPath, "config.json")

	if err := os.WriteFile(configPath, jsonData, 0640); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	return nil
}

func (f *filesystemProblemCodeRepository) GetProblemCodes(ctx context.Context, problem *models.Problem) ([]*models.ProblemCode, error) {
	problemCodes := make([]*models.ProblemCode, 0, models.LanguageEnd-1)
	for lang := models.LanguageGo; lang < models.LanguageEnd; lang++ {
		problemCode := new(models.ProblemCode)
		problemCode.Language = lang

		var err error
		problemCode.Checker, err = f.getLanguageFile(problem, lang, lang.CheckerFilename())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		problemCode.CodeSnippet, err = f.getLanguageFile(problem, lang, lang.CodeSnippetFilename())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		problemCode.Driver, err = f.getLanguageFile(problem, lang, lang.DriverFilename())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		problemCodeMetadata, err := f.getMetadataLanguageFile(problem, lang)
		if err != nil {
			return nil, err
		}
		problemCode.IsPublic = problemCodeMetadata.IsPublic

		problemCodes = append(problemCodes, problemCode)
	}

	return problemCodes, nil
}

func (f *filesystemProblemCodeRepository) SaveProblemCode(ctx context.Context, problem *models.Problem, problemCode *models.ProblemCode) error {
	err := f.saveLanguageFile(
		problem,
		problemCode.Language,
		problemCode.Language.CodeSnippetFilename(),
		[]byte(problemCode.CodeSnippet),
	)
	if err != nil {
		return err
	}

	err = f.saveLanguageFile(
		problem,
		problemCode.Language,
		problemCode.Language.CheckerFilename(),
		[]byte(problemCode.Checker),
	)
	if err != nil {
		return err
	}

	err = f.saveLanguageFile(
		problem,
		problemCode.Language,
		problemCode.Language.DriverFilename(),
		[]byte(problemCode.Driver),
	)
	if err != nil {
		return err
	}

	err = f.saveJsonLanguageFile(
		problem,
		problemCode.Language,
		"metadata.json",
		&models.ProblemCodeMetadata{
			IsPublic: problemCode.IsPublic,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (f *filesystemProblemCodeRepository) getLanguageFile(
	problem *models.Problem,
	language models.ProblemLanguage,
	filename string,
) (string, error) {
	languagePath, err := f.getLanguagePath(
		problem,
		language,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get language path: %w", err)
	}
	filePath := filepath.Join(languagePath, filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (f *filesystemProblemCodeRepository) saveLanguageFile(
	problem *models.Problem,
	language models.ProblemLanguage,
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

func (f *filesystemProblemCodeRepository) saveJsonLanguageFile(
	problem *models.Problem,
	language models.ProblemLanguage,
	filename string,
	content any,
) error {
	jsonData, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal indent content: %w", err)
	}

	languagePath, err := f.getLanguagePath(problem, language)
	if err != nil {
		return fmt.Errorf("failed to get languagePath: %w", err)
	}
	filePath := filepath.Join(languagePath, filename)

	if err := os.WriteFile(filePath, jsonData, 0640); err != nil {
		return fmt.Errorf("failed to write json to file: %w", err)
	}

	return nil
}

func (f *filesystemProblemCodeRepository) getMetadataLanguageFile(
	problem *models.Problem,
	language models.ProblemLanguage,
) (*models.ProblemCodeMetadata, error) {
	languagePath, err := f.getLanguagePath(problem, language)
	if err != nil {
		return nil, fmt.Errorf("failed to get languagePath: %w", err)
	}
	metadataPath := filepath.Join(languagePath, "metadata.json")

	file, err := os.Open(metadataPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			metadata := models.DefaultProblemCodeMetadata()
			if err := f.saveJsonLanguageFile(problem, language, "metadata.json", metadata); err != nil {
				return nil, fmt.Errorf("failed to save default metadata: %w", err)
			}

			return metadata, nil
		}

		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var metadata *models.ProblemCodeMetadata
	err = decoder.Decode(&metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata file: %w", err)
	}

	return metadata, nil
}

func (f *filesystemProblemCodeRepository) SaveTestcases(
	ctx context.Context,
	problem *models.Problem,
	problemTestcases []*models.ProblemTestcase,
) error {
	jsonData, err := json.MarshalIndent(problemTestcases, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal indent problemTestcases: %w", err)
	}

	problemPath, err := f.getProblemPath(problem)
	if err != nil {
		return fmt.Errorf("failed to get problemPath: %w", err)
	}
	testcasesPath := filepath.Join(problemPath, "testcases.json")

	if err := os.WriteFile(testcasesPath, jsonData, 0640); err != nil {
		return fmt.Errorf("failed to write testcases to file: %w", err)
	}

	return nil
}

func (f *filesystemProblemCodeRepository) GetTestcases(
	ctx context.Context,
	problem *models.Problem,
) ([]*models.ProblemTestcase, error) {
	problemPath, err := f.getProblemPath(problem)
	if err != nil {
		return nil, fmt.Errorf("failed to get problemPath: %w", err)
	}
	testcasesPath := filepath.Join(problemPath, "testcases.json")

	file, err := os.Open(testcasesPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			testcases := []*models.ProblemTestcase{}
			if err := f.SaveTestcases(ctx, problem, testcases); err != nil {
				return nil, fmt.Errorf("failed to save default testcases: %w", err)
			}

			return testcases, nil
		}

		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var testcases []*models.ProblemTestcase
	err = decoder.Decode(&testcases)
	if err != nil {
		return nil, fmt.Errorf("failed to decode testcases file: %w", err)
	}

	return testcases, nil
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

func (f *filesystemProblemCodeRepository) getLanguagePath(problem *models.Problem, language models.ProblemLanguage) (string, error) {
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
