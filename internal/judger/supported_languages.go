package judger

import (
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

type Language struct {
	containerImage string
	internalSlug   string
	new            func(cli *client.Client) languageJudge
	logger         *slog.Logger
}

var problemsDataDir string

func init() {
	problemsDataDir = os.Getenv("CODEUNIVERSE_PROBLEMS_DATA_DIR")

	if problemsDataDir == "" {
		log.Fatal("CODEUNIVERSE_PROBLEMS_DATA_DIR is not set.")
	}

	absPath, err := filepath.Abs(problemsDataDir)
	if err != nil {
		log.Fatal("failed to convert CODEUNIVERSE_PROBLEMS_DATA_DIR to absolute path.")
	}

	problemsDataDir = absPath
	slog.Info("problemsDataDir is updated.", "problemsDataDir", problemsDataDir)
}

var SupportedLanguages = map[string]Language{
	"c++": {
		"gcc:15.2-trixie",
		"cpp",
		newCPPJudge,
		slog.Default().With("package", "judge.cppLanguage")},
	"golang": {
		"golang:1.25.4-alpine",
		"go",
		newGolangJudge,
		slog.Default().With("package", "judge.goLanguage")},
	"python3": {
		"python:3.13.7-alpine",
		"py3",
		newPython3Judge,
		slog.Default().With("package", "judge.py3Language"),
	},
	"javascript": {
		"node:22.21.1-alpine",
		"js",
		newJavascriptJudge,
		slog.Default().With("package", "judge.jsLanguage"),
	},
	"typescript": {
		"node:22.21.1-alpine",
		"ts",
		newTypescriptJudge,
		slog.Default().With("package", "judge.tsLanguage"),
	},
}

func (l *Language) DoesItHaveTests(problemSlug string) bool {
	workspaceDir := filepath.Join(problemsDataDir, "problems", problemSlug, l.internalSlug)

	object, err := os.Stat(workspaceDir)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return false
		default:
			l.logger.Error(
				"failed to see stats for workspaceDir",
				"workspaceDir",
				workspaceDir,
				"err",
				err,
			)
		}
	}

	isDir := object.IsDir()
	if !isDir {
		l.logger.Error(
			"seems that we have a file instead of a directory",
			"workspaceDir",
			workspaceDir,
		)
	}

	return isDir
}
