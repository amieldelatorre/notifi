package middleware

import (
	"io"
	"log/slog"

	"github.com/amieldelatorre/notifi/logger"
)

func GetMockMiddleware() Middleware {
	logger := logger.New(io.Discard, slog.LevelWarn)
	return Middleware{Logger: logger}
}
