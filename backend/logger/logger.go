package logger

import (
	"io"
	"log/slog"
)

func New(w io.Writer, level slog.Leveler) *slog.Logger {
	opts := slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(w, &opts))
	return logger
}
