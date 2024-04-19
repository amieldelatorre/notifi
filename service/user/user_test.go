package user

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/service/security"
	"github.com/jackc/pgx/v5"
)

func TestGetUserById(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	mockUserProvider := NewMockUserRepo()
	service := New(logger, &mockUserProvider)

	testCases := GetValidTestGetUserByIdTestCases()
	testCases = append(testCases, GetInvalidTestGetUserByIdTestCase())

	for _, tc := range testCases {
		actualStatusCode, actualResponse := service.GetUserById(context.Background(), tc.UserId)

		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case userId %d, expected exected status code %d, got %d", tc.UserId, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.Response.Errors) {
			t.Fatalf("test case userId %d, expected response errors %+v, got %+v", tc.UserId, tc.Response.Errors, actualResponse.Errors)
		}

		if actualResponse.User != nil && tc.Response.User != nil && (actualResponse.User.Id != tc.Response.User.Id || actualResponse.User.Email != tc.Response.User.Email ||
			actualResponse.User.FirstName != tc.Response.User.FirstName || actualResponse.User.LastName != tc.Response.User.LastName ||
			actualResponse.User.Password != tc.Response.User.Password) {
			t.Fatalf("test case userId %d, expected response user %+v, got %+v", tc.UserId, tc.Response.User, actualResponse.User)
		}
	}
}

func TestGetUserByEmail(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	mockUserProvider := NewMockUserRepo()
	service := New(logger, &mockUserProvider)

	tcs := []struct {
		Email          string
		ExpectedError  error
		ExpectedResult model.User
	}{
		{
			Email:          mockUserProvider.Users[1].Email,
			ExpectedError:  nil,
			ExpectedResult: mockUserProvider.Users[1],
		},
		{
			Email:          "email@notexist.invalid",
			ExpectedError:  pgx.ErrNoRows,
			ExpectedResult: model.User{},
		},
	}

	for _, tc := range tcs {
		actualResult, actualError := service.Provider.GetUserByEmail(context.Background(), tc.Email)
		if actualResult != tc.ExpectedResult {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.Email, tc.ExpectedResult, actualResult)
		}

		if actualError != tc.ExpectedError {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.Email, tc.ExpectedError, actualError)
		}
	}
}

func TestEmailExists(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	mockUserProvider := NewMockUserRepo()
	service := New(logger, &mockUserProvider)

	tcs := []struct {
		Email          string
		ExpectedError  error
		ExpectedResult bool
	}{
		{
			Email:          mockUserProvider.Users[1].Email,
			ExpectedError:  nil,
			ExpectedResult: true,
		},
		{
			Email:          "email@notexist.invalid",
			ExpectedError:  nil,
			ExpectedResult: false,
		},
	}

	for _, tc := range tcs {
		actualResult, actualError := service.EmailExists(context.Background(), tc.Email)
		if actualResult != tc.ExpectedResult {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.Email, tc.ExpectedResult, actualResult)
		}

		if actualError != tc.ExpectedError {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.Email, tc.ExpectedError, actualError)
		}
	}
}

func TestValidateUserInput(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	mockUserProvider := NewMockUserRepo()
	service := New(logger, &mockUserProvider)

	tcs := []struct {
		UserInput                model.UserInput
		ExpectedError            error
		ExpectedValidationErrors map[string][]string
		ExpectedCleanUserInput   model.UserInput
	}{
		{
			UserInput: model.UserInput{
				Email:     mockUserProvider.Users[0].Email,
				FirstName: mockUserProvider.Users[0].FirstName,
				LastName:  mockUserProvider.Users[0].LastName,
				Password:  mockUserProvider.Users[0].Password,
			},
			ExpectedError:            nil,
			ExpectedValidationErrors: map[string][]string{"email": {"Email already exists"}},
			ExpectedCleanUserInput: model.UserInput{
				Email:     mockUserProvider.Users[0].Email,
				FirstName: mockUserProvider.Users[0].FirstName,
				LastName:  mockUserProvider.Users[0].LastName,
				Password:  mockUserProvider.Users[0].Password,
			},
		},
		{
			UserInput: model.UserInput{
				Email:     "email@notexist.invalid",
				FirstName: mockUserProvider.Users[0].FirstName,
				LastName:  mockUserProvider.Users[0].LastName,
				Password:  mockUserProvider.Users[0].Password,
			},
			ExpectedError:            nil,
			ExpectedValidationErrors: map[string][]string{},
			ExpectedCleanUserInput: model.UserInput{
				Email:     "email@notexist.invalid",
				FirstName: mockUserProvider.Users[0].FirstName,
				LastName:  mockUserProvider.Users[0].LastName,
				Password:  mockUserProvider.Users[0].Password,
			},
		},
	}

	for _, tc := range tcs {
		actualUserInput, actualValidationErrors, actualError := service.validateUserinput(context.Background(), tc.UserInput)
		if actualUserInput != tc.ExpectedCleanUserInput {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.UserInput.Email, tc.ExpectedCleanUserInput, actualUserInput)
		}

		if !reflect.DeepEqual(actualValidationErrors, tc.ExpectedValidationErrors) {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.UserInput.Email, tc.ExpectedValidationErrors, actualValidationErrors)
		}

		if actualError != tc.ExpectedError {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.UserInput.Email, tc.ExpectedError, actualError)
		}
	}
}

func TestCreateUserInput(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	mockUserProvider := NewMockUserRepo()
	service := New(logger, &mockUserProvider)

	tcs := []struct {
		UserInput          model.UserInput
		ExpectedStatusCode int
		ExpectedResponse   UserResponse
	}{
		{
			UserInput: model.UserInput{
				Email:     mockUserProvider.Users[0].Email,
				FirstName: mockUserProvider.Users[0].FirstName,
				LastName:  mockUserProvider.Users[0].LastName,
				Password:  mockUserProvider.Users[0].Password,
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: UserResponse{
				Errors: map[string][]string{"email": {"Email already exists"}},
			},
		},
		{
			UserInput: model.UserInput{
				Email:     "",
				FirstName: "",
				LastName:  "",
				Password:  "",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: UserResponse{
				Errors: map[string][]string{
					"email":     {"Cannot be empty"},
					"firstName": {"Cannot be empty"},
					"lastName":  {"Cannot be empty"},
					"password":  {"Cannot be empty and must be at least 8 characters"},
				},
			},
		},
		{
			UserInput: model.UserInput{
				Email:     "isaac.newton@invalid.com",
				FirstName: "Isaac      ",
				LastName:  "Newton",
				Password:  "Password",
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse: UserResponse{
				Errors: map[string][]string{},
				User: &model.User{
					Email:           "isaac.newton@invalid.com",
					Id:              len(mockUserProvider.Users) + 1,
					FirstName:       "Isaac",
					LastName:        "Newton",
					Password:        "Password",
					DatetimeCreated: time.Now(),
					DatetimeUpdated: time.Now(),
				},
			},
		},
	}

	for _, tc := range tcs {
		actualStatusCode, actualResponse := service.CreateUser(context.Background(), tc.UserInput)
		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case userId %s, expected response user %d, got %d", tc.UserInput.Email, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.UserInput.Email, tc.ExpectedResponse, actualResponse)
		}

		if actualResponse.User != nil && tc.ExpectedResponse.User != nil && (actualResponse.User.Id != tc.ExpectedResponse.User.Id || actualResponse.User.Email != tc.ExpectedResponse.User.Email ||
			actualResponse.User.FirstName != tc.ExpectedResponse.User.FirstName || actualResponse.User.LastName != tc.ExpectedResponse.User.LastName) {
			t.Fatalf("test case userId %s, expected response user %+v, got %+v", tc.UserInput.Email, tc.ExpectedResponse.User, actualResponse.User)
		}

		if actualResponse.User != nil && tc.ExpectedResponse.User != nil {
			passwordsMatch, err := security.IsCorrectPassword(context.Background(), tc.UserInput.Password, actualResponse.User.Password, logger)
			if err != nil {
				t.Error(err)
			}

			if !passwordsMatch {
				t.Fatalf("expected passwords to match")
			}
		}
	}
}
