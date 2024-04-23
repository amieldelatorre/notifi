package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/service/security"
	"github.com/amieldelatorre/notifi/utils"
)

func GetMockMiddleware() Middleware {
	logger := logger.New(io.Discard, slog.LevelWarn)
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))
	return Middleware{Logger: logger, JwtService: &jwtService}
}

func TestRecoverPanic(t *testing.T) {
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Test something going wrong")
	})

	m := GetMockMiddleware()
	handler := m.RecoverPanic(mockNextHandler)
	request := httptest.NewRequest("GET", "http://test.invalid", nil)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	result := recorder.Result()

	expectedStatusCode := http.StatusInternalServerError

	if result.StatusCode != expectedStatusCode {
		t.Fatalf("expected status code '%d', got '%d'", expectedStatusCode, result.StatusCode)
	}
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
	request := httptest.NewRequest("GET", "http://test.invalid", nil)

	handler.ServeHTTP(httptest.NewRecorder(), request)

}

func TestRequireJwt(t *testing.T) {
	m := GetMockMiddleware()
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tcs := []struct {
		ExpectedStatusCode int
		ExpectedMessage    string
		HeaderValue        *string
	}{}

	// Test Case 1
	claims1 := security.UserClaims{
		UserId: 1,
		Email:  "email@test.invalid",
	}
	token1, err := m.JwtService.CreateAccessToken(claims1)
	if err != nil {
		t.Error(err)
	}
	headerValue1 := fmt.Sprintf("Bearer %s", token1)

	tcs = append(tcs, struct {
		ExpectedStatusCode int
		ExpectedMessage    string
		HeaderValue        *string
	}{
		ExpectedStatusCode: http.StatusOK,
		HeaderValue:        &headerValue1,
	})

	// Test Case 2
	tcs = append(tcs, struct {
		ExpectedStatusCode int
		ExpectedMessage    string
		HeaderValue        *string
	}{
		ExpectedStatusCode: http.StatusUnauthorized,
		ExpectedMessage:    "Missing Authorization header",
	})

	// Test Case 3
	headerValue3 := "Bearer"
	tcs = append(tcs, struct {
		ExpectedStatusCode int
		ExpectedMessage    string
		HeaderValue        *string
	}{
		ExpectedStatusCode: http.StatusUnauthorized,
		ExpectedMessage:    "Wrong Authorization header type, this endpoint requires an Authorization of 'Bearer' and a token",
		HeaderValue:        &headerValue3,
	})

	// Test Case 4
	headerValue4 := "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	tcs = append(tcs, struct {
		ExpectedStatusCode int
		ExpectedMessage    string
		HeaderValue        *string
	}{
		ExpectedStatusCode: http.StatusUnauthorized,
		ExpectedMessage:    "Wrong Authorization header type, this endpoint requires an Authorization of 'Bearer'",
		HeaderValue:        &headerValue4,
	})

	// Test Case 5
	headerValue5 := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	tcs = append(tcs, struct {
		ExpectedStatusCode int
		ExpectedMessage    string
		HeaderValue        *string
	}{
		ExpectedStatusCode: http.StatusUnauthorized,
		ExpectedMessage:    "Invalid 'Bearer' token",
		HeaderValue:        &headerValue5,
	})

	for tcn, tc := range tcs {
		handler := m.RequireJwtToken(mockNextHandler)
		request := httptest.NewRequest("GET", "http://test.invalid", nil)
		if tc.HeaderValue != nil {
			request.Header.Add("Authorization", *tc.HeaderValue)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		result := recorder.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case %d, expected status code '%d', got '%d'", tcn, tc.ExpectedStatusCode, result.StatusCode)
		}

		resultBody, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		result.Body.Close()

		if tc.ExpectedStatusCode != http.StatusOK {
			var jwtErrors RequireJwtTokenErrors
			err = json.Unmarshal(resultBody, &jwtErrors)
			if err != nil {
				t.Error(err)
			}

			if tc.ExpectedMessage != jwtErrors.Errors["Authorization"][0] {
				t.Fatalf("test case %d, expected status code '%s', got '%s'", tcn, tc.ExpectedMessage, jwtErrors.Errors["Authorization"][0])
			}
		}

	}
}
