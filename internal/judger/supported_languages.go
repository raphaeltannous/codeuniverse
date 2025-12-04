package judger

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
)

type Language struct {
	containerImage string
	internalSlug   string
	new            func(cli *client.Client) languageJudge
	logger         *slog.Logger
}

const problemsDataDir = "problems-data"

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
	workspaceDir := filepath.Join(problemsDataDir, problemSlug, l.internalSlug)

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

func (l *Language) copyRunToWorkspace(problemSlug, workspace string) error {
	return l.copyToWorkspace(problemSlug, workspace, "run")
}

func (l *Language) copySubmitToWorkspace(problemSlug, workspace string) error {
	return l.copyToWorkspace(problemSlug, workspace, "submit")
}

func (l *Language) copyToWorkspace(problemSlug, workspace, executionType string) error {
	problemTestDir := filepath.Join(problemsDataDir, problemSlug, l.internalSlug)
	srcDir := os.DirFS(problemTestDir)

	err := os.CopyFS(workspace, srcDir)
	if err != nil {
		return err
	}

	return filepath.WalkDir(workspace, func(path string, de fs.DirEntry, err error) error {
		if !de.IsDir() {
			filename := de.Name()

			if strings.HasSuffix(filename, "."+executionType+".tmpl") {
				if err := os.Rename(path, path[:len(path)-len("."+executionType+".tmpl")]); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
