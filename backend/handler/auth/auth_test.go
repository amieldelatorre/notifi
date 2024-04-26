package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amieldelatorre/notifi/backend/logger"
	authService "github.com/amieldelatorre/notifi/backend/service/auth"
	"github.com/amieldelatorre/notifi/backend/service/security"
	userService "github.com/amieldelatorre/notifi/backend/service/user"
)

func GetNewMockAuthHandler() (AuthHandler, userService.TestDbProviderInstance) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := userService.NewTestDbInstance()
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))

	authService := authService.New(logger, testDbInstance.Provider, jwtService)

	mockAuthHandler := New(logger, authService, jwtService)
	return mockAuthHandler, testDbInstance
}

func TestLogin(t *testing.T) {
	mockAuthHandler, testDbInstance := GetNewMockAuthHandler()
	defer testDbInstance.CleanUp()

	tcs := []struct {
		Body struct {
			Email    string
			Password string
		}
		ExpectedStatusCode int
	}{
		{
			Body: struct {
				Email    string
				Password string
			}{
				Email:    "isaac.newton@example.invalid",
				Password: "password",
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Body: struct {
				Email    string
				Password string
			}{
				Email:    "email.notexists@example.invalid",
				Password: "password",
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tcs {
		body, err := json.Marshal(tc.Body)
		if err != nil {
			t.Error(err)
		}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		mockAuthHandler.login(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case %s, expected status code '%d', got '%d'", tc.Body.Email, tc.ExpectedStatusCode, result.StatusCode)
		}
	}
}
