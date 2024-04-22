package middleware // import "github.com/amieldelatorre/notifi/middleware"

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/utils"
	"github.com/google/uuid"
)

type Middleware struct {
	Logger *slog.Logger
}

func New(logger *slog.Logger) Middleware {
	return Middleware{Logger: logger}
}

func (m *Middleware) AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Logger.Debug("Adding a request id to incoming request")
		id, err := uuid.NewV7()
		if err != nil {
			m.Logger.Error("Problem generating uuid: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.Logger.Debug("Request id generated and added to context", "requestId", id.String())

		ctx := context.WithValue(r.Context(), utils.RequestIdName, id.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) RequireJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Context().Value(utils.RequestIdName)
		m.Logger.Debug("Checking if request has the Authorization Bearer token", "requestId", requestId)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Validate token
		// If not valid, return Unauthorized
		// If valid set the UserId in context
		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
