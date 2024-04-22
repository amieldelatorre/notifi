package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amieldelatorre/notifi/logger"
)

func GetMockMiddleware() Middleware {
	logger := logger.New(io.Discard, slog.LevelWarn)
	return Middleware{Logger: logger}
}

func TestAddRequestId(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(RequestIdName)
		if val == nil {
			t.Fatalf("'%s' is missing", RequestIdName)
		}

		value, ok := val.(string)
		if !ok {
			t.Fatalf("Value of '%s', %+v, is not a string", RequestIdName, value)
		}
	})

	m := GetMockMiddleware()
	handler := m.AddRequestId(nextHandler)
	request := httptest.NewRequest("GET", "http://example.invalid", nil)

	handler.ServeHTTP(httptest.NewRecorder(), request)

}
