package middleware

import "github.com/amieldelatorre/notifi/logger"

type Middleware struct {
	Logger *logger.Logger
}

func New(logger *logger.Logger) Middleware {
	return Middleware{Logger: logger}
}
