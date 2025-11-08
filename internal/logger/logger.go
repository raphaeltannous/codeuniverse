package logger

import (
	"io"
	"log/slog"
	"os"
)

func New(level slog.Leveler) (*slog.Logger, error) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(
		io.MultiWriter(os.Stdout, logFile),
		&slog.HandlerOptions{Level: level},
	)

	return slog.New(handler), nil
}
