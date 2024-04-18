package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
