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
	// The structure of the baseDirectory will be as follow:
	// - baseDirectory/problems/problemUUID/					Contains user codesnippet and wanted files for docker.
	//   - /languageSlug										Seperate each language alone.
	//     - /language.DriverFilename							The file that loads the testcases.
	//     - /language.CheckerFilename							Code that accepts testcases and returns result.
	//     - /language.CodeSnippetFilename						User facing code snippet for the language.
	//     - /metadata.json										Currently contains if the language is public or not.
	//   - /config.json											Docker container config. (TimeLimit/MemoryLimit).
	//   - /testcases.json
	// - baseDirectory/languages/
	//   - /templates/languageSlug								Templates for the creating a new problem.
	//     - /language.DriverFilename
	//     - /language.CheckerFilename
	//     - /language.CodeSnippetFilename
	//   - /runtime/languageSlug								Files that are copied to the docker container and are available at /codeuniverse/run.
	//     - /language.BackendCheckerFilename					The file that will import and call the function in language.CheckerFilename file.
	//     - other libs.

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
		problemCode, err := f.GetProblemCode(ctx, problem, lang)
		if err != nil {
			return nil, err
		}

		problemCodes = append(problemCodes, problemCode)
	}

	return problemCodes, nil
}

func (f *filesystemProblemCodeRepository) getRuntimeFilesList(
	language models.ProblemLanguage,
) ([]string, error) {
	runtimeLanguageDir := filepath.Join(
		f.baseDirectory,
		"languages",
		"runtime",
		language.Slug(),
	)

	dirStat, err := os.Stat(runtimeLanguageDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}

		return []string(nil), nil
	}

	if dirStat.IsDir() {
		files, err := os.ReadDir(runtimeLanguageDir)
		if err != nil {
			return []string(nil), nil
		}

		runtimeFiles := make([]string, 0, len(files))
		for _, file := range files {
			if !file.IsDir() {
				runtimeFiles = append(
					runtimeFiles,
					filepath.Join(runtimeLanguageDir, file.Name()),
				)
			}
		}

		return runtimeFiles, nil
	}

	return []string(nil), nil
}

func (f *filesystemProblemCodeRepository) GetProblemCode(ctx context.Context, problem *models.Problem, language models.ProblemLanguage) (*models.ProblemCode, error) {
	problemCode := new(models.ProblemCode)
	problemCode.Language = language

	var err error
	problemCode.Checker, err = f.getLanguageFile(problem, language, language.CheckerFilename())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		problemCode.Checker, err = f.getLanguageTemplateFile(language, language.CheckerFilename())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	problemCode.CodeSnippet, err = f.getLanguageFile(problem, language, language.CodeSnippetFilename())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		problemCode.CodeSnippet, err = f.getLanguageTemplateFile(language, language.CodeSnippetFilename())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	problemCode.Driver, err = f.getLanguageFile(problem, language, language.DriverFilename())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		problemCode.Driver, err = f.getLanguageTemplateFile(language, language.DriverFilename())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	problemCode.RuntimeFiles, err = f.getRuntimeFilesList(language)
	if err != nil {
		return nil, err
	}

	problemCodeMetadata, err := f.getMetadataLanguageFile(problem, language)
	if err != nil {
		return nil, err
	}
	problemCode.IsPublic = problemCodeMetadata.IsPublic

	return problemCode, nil
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

func (f *filesystemProblemCodeRepository) getLanguageTemplateFile(
	language models.ProblemLanguage,
	filename string,
) (string, error) {
	languageTemplatePath, err := f.getLanguageTemplatePath(
		language,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get language template path: %w", err)
	}
	filePath := filepath.Join(languageTemplatePath, filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
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
	problemPath := filepath.Join(f.baseDirectory, "problems", problem.ID.String())

	if _, err := os.Stat(problemPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(problemPath, 0750); err != nil {
			return "", err
		}
	}

	return problemPath, nil
}

func (f *filesystemProblemCodeRepository) getLanguageTemplatePath(language models.ProblemLanguage) (string, error) {
	templatesPath := filepath.Join(
		f.baseDirectory,
		"languages",
		"templates",
	)

	languageTemplatePath := filepath.Join(templatesPath, language.Slug())
	if _, err := os.Stat(languageTemplatePath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(languageTemplatePath, 0750); err != nil {
			return "", err
		}
	}

	return languageTemplatePath, nil
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
