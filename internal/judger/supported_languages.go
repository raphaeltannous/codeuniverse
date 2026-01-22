package judger

import (
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/docker/docker/client"
)

type LanguageJudge struct {
	containerImage string
	judge          languageJudger
	logger         *slog.Logger
}

var supportedLanguagesNewFunc = map[models.ProblemLanguage]func(cli *client.Client, logger *slog.Logger) languageJudger{
	models.LanguageGo:         newGolangJudge,
	models.LanguagePython:     newPythonJudge,
	models.LanguageCpp:        newCPPJudge,
	models.LanguageTypescript: newTypescriptJudge,
	models.LanguageJavascript: newJavascriptJudge,
	models.LanguageJava:       newJavaJudge,
	models.LanguageRuby:       newRubyJudge,
}

var supportedLanguages = map[models.ProblemLanguage]*LanguageJudge{
	models.LanguageGo: {
		containerImage: "golang:1.25.4-alpine",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.goLanguage"),
	},
	models.LanguagePython: {
		containerImage: "python:3.13.7-alpine",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.pythonLanguage"),
	},
	models.LanguageCpp: {
		containerImage: "gcc:15.2-trixie",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.cppLanguage"),
	},
	models.LanguageTypescript: {
		containerImage: "codeuniverse-node:22.21.1-alpine",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.typescriptLanguage"),
	},
	models.LanguageJavascript: {
		containerImage: "codeuniverse-node:22.21.1-alpine",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.javascriptLanguage"),
	},
	models.LanguageJava: {
		containerImage: "amazoncorretto:21-alpine-jdk",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.javaLanguage"),
	},
	models.LanguageRuby: {
		containerImage: "ruby:3.4.8-alpine",
		judge:          nil,
		logger:         slog.Default().With("package", "judger.rubyLanguage"),
	},
}
