package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

var LoggerWriter io.Writer

func init() {
	// TODO: change app.log to its correct path.
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		log.Fatalf("failed to initialize logFile: %v", err)
	}

	LoggerWriter = io.MultiWriter(os.Stdout, logFile)
}

func New(level slog.Leveler) (*slog.Logger, error) {
	handler := slog.NewJSONHandler(
		LoggerWriter,
		&slog.HandlerOptions{Level: level},
	)

	return slog.New(handler), nil
}
