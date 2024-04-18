package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"time"
)

type Logger struct {
	Logger *slog.Logger
}

func New(w io.Writer, level slog.Leveler) *Logger {
	opts := slog.HandlerOptions{
		Level: level,
	}

	slogLogger := slog.New(slog.NewJSONHandler(w, &opts))

	return &Logger{Logger: slogLogger}
}

func (l *Logger) Log(level slog.Level, format string, args ...any) {
	if !l.Logger.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, log]

	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
	_ = l.Logger.Handler().Handle(context.Background(), r)
}

func (l *Logger) Info(format string, args ...any) {
	l.Log(slog.LevelInfo, format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	l.Log(slog.LevelDebug, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.Log(slog.LevelWarn, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.Log(slog.LevelError, format, args...)
}
