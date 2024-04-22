package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/utils"
)

func GetMockMiddleware() Middleware {
	logger := logger.New(io.Discard, slog.LevelWarn)
	return Middleware{Logger: logger}
}

func TestAddRequestId(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(utils.RequestIdName)
		if val == nil {
			t.Fatalf("'%s' is missing", utils.RequestIdName)
		}

		value, ok := val.(string)
		if !ok {
			t.Fatalf("Value of '%s', %+v, is not a string", utils.RequestIdName, value)
		}
	})

	m := GetMockMiddleware()
	handler := m.AddRequestId(nextHandler)
	request := httptest.NewRequest("GET", "http://example.invalid", nil)

	handler.ServeHTTP(httptest.NewRecorder(), request)

}
