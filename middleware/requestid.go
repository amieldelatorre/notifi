package middleware // import "github.com/amieldelatorre/notifi/middleware"

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

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

		ctx := context.WithValue(r.Context(), RequestIdName, id.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
