package middleware // import "github.com/amieldelatorre/notifi/middleware"

import (
	"log/slog"
)

type Middleware struct {
	Logger *slog.Logger
}

func New(logger *slog.Logger) Middleware {
	return Middleware{Logger: logger}
}
