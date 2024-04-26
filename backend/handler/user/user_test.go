package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/backend/logger"
	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/amieldelatorre/notifi/backend/service/security"
	userService "github.com/amieldelatorre/notifi/backend/service/user"
	"github.com/amieldelatorre/notifi/backend/utils"
)

func GetNewMockUserHandler() (UserHandler, userService.TestDbProviderInstance) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := userService.NewTestDbInstance()
	usrService := userService.New(logger, testDbInstance.Provider)
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))

	mockUserHandler := New(logger, usrService, jwtService)
	return mockUserHandler, testDbInstance
}

func TestGetUser(t *testing.T) {
	mockUserHandler, testDbInstance := GetNewMockUserHandler()
	defer testDbInstance.CleanUp()

	testcases := userService.GetValidTestGetUserByIdTestCases()
	testcases = append(testcases, userService.GetInvalidTestGetUserByIdTestCase())

	for _, tc := range testcases {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, utils.UserId, tc.UserId)
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		mockUserHandler.getUser(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case userId %d, expected status code %d, got %d", tc.UserId, tc.ExpectedStatusCode, result.StatusCode)
		}

		body, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		result.Body.Close()

		var userResponse userService.UserResponse
		err = json.Unmarshal(body, &userResponse)
		if err != nil {
			t.Error(err)
		}

		// Check if both are length of 0 as userHandler.getUser does omits the Errors if it is empty
		if len(userResponse.Errors) != 0 && len(tc.Response.Errors) != 0 && !reflect.DeepEqual(userResponse.Errors, tc.Response.Errors) {
			t.Fatalf("test case userId %d, expected response errors %+v, got %+v", tc.UserId, tc.Response.Errors, userResponse.Errors)
		}

		if userResponse.User != nil && tc.Response.User != nil && (userResponse.User.Id != tc.Response.User.Id || userResponse.User.Email != tc.Response.User.Email ||
			userResponse.User.FirstName != tc.Response.User.FirstName || userResponse.User.LastName != tc.Response.User.LastName) {

			jsonExpectedUser, err := json.Marshal(tc.Response.User)
			if err != nil {
				t.Error(err)
			}

			jsonResponseUser, err := json.Marshal(userResponse.User)
			if err != nil {
				t.Error(err)
			}
			t.Fatalf("test case userId %d, expected response user %+v, got %+v", tc.UserId, string(jsonExpectedUser), string(jsonResponseUser))
		}
	}
}

func TestPostUser(t *testing.T) {
	mockUserHandler, testDbInstance := GetNewMockUserHandler()
	defer testDbInstance.CleanUp()

	tcs := []struct {
		UserInput          model.UserInput
		ExpectedStatusCode int
		ExpectedResponse   userService.UserResponse
	}{
		{
			UserInput: model.UserInput{
				Email:     "isaac.newton@example.invalid",
				FirstName: "Isaac",
				LastName:  "Newton",
				Password:  "Password",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: userService.UserResponse{
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
			ExpectedResponse: userService.UserResponse{
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
			ExpectedResponse: userService.UserResponse{
				Errors: map[string][]string{},
				User: &model.User{
					Email:           "isaac.newton@invalid.com",
					Id:              3,
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
		body, err := json.Marshal(tc.UserInput)
		if err != nil {
			t.Error(err)
		}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		mockUserHandler.postUser(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case userEmail %s, expected status code %d, got %d", tc.UserInput.Email, tc.ExpectedStatusCode, result.StatusCode)
		}

		resultBody, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		result.Body.Close()

		var userResponse userService.UserResponse
		err = json.Unmarshal(resultBody, &userResponse)
		if err != nil {
			t.Error(err)
		}

		// Check if both are length of 0 as userHandler.getUser does omits the Errors if it is empty
		if len(userResponse.Errors) != 0 && len(tc.ExpectedResponse.Errors) != 0 && !reflect.DeepEqual(userResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case userEmail %s, expected response errors %+v, got %+v", tc.UserInput.Email, tc.ExpectedResponse.Errors, userResponse.Errors)
		}

		if userResponse.User != nil && tc.ExpectedResponse.User != nil && (userResponse.User.Id != tc.ExpectedResponse.User.Id || userResponse.User.Email != tc.ExpectedResponse.User.Email ||
			userResponse.User.FirstName != tc.ExpectedResponse.User.FirstName || userResponse.User.LastName != tc.ExpectedResponse.User.LastName) {

			jsonExpectedUser, err := json.Marshal(tc.ExpectedResponse.User)
			if err != nil {
				t.Error(err)
			}

			jsonResponseUser, err := json.Marshal(userResponse.User)
			if err != nil {
				t.Error(err)
			}
			t.Fatalf("test case userEmail %s, expected response user %+v, got %+v", tc.UserInput.Email, string(jsonExpectedUser), string(jsonResponseUser))
		}
	}
}
