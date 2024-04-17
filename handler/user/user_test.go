package user

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	userService "github.com/amieldelatorre/notifi/service/user"
)

func GetNewMockUserHandler() UserHandler {
	mockUserProvider := userService.NewMockUserRepo()
	service := userService.New(&mockUserProvider)
	mockUserHandler := New(service)
	return mockUserHandler
}

func TestGetUserMissingUserIdHeader(t *testing.T) {
	mockUserHandler := GetNewMockUserHandler()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
	response := httptest.NewRecorder()
	expectedStatusCode := http.StatusInternalServerError

	mockUserHandler.getUser(response, request)

	result := response.Result()
	if result.StatusCode != expectedStatusCode {
		t.Fatalf("expected status code %d, got %d", expectedStatusCode, result.StatusCode)
	}
}

func TestGetUser(t *testing.T) {
	mockUserHandler := GetNewMockUserHandler()

	testcases := userService.GetValidTestGetUserByIdTestCases()
	testcases = append(testcases, userService.GetInvalidTestGetUserByIdTestCase())

	for _, tc := range testcases {
		userId := strconv.Itoa(tc.UserId)

		request := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
		request.Header.Set("x-user-id", userId)

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

		var getUserResponse userService.GetUserResponse
		err = json.Unmarshal(body, &getUserResponse)
		if err != nil {
			t.Error(err)
		}

		// Check if both are length of 0 as userHandler.getUser does omits the Errors if it is empty
		if len(getUserResponse.Errors) != 0 && len(tc.Response.Errors) != 0 && !reflect.DeepEqual(getUserResponse.Errors, tc.Response.Errors) {
			t.Fatalf("test case userId %d, expected response errors %+v, got %+v", tc.UserId, tc.Response.Errors, getUserResponse.Errors)
		}

		if getUserResponse.User != nil && tc.Response.User != nil && (getUserResponse.User.Id != tc.Response.User.Id || getUserResponse.User.Email != tc.Response.User.Email ||
			getUserResponse.User.FirstName != tc.Response.User.FirstName || getUserResponse.User.LastName != tc.Response.User.LastName) {

			jsonExpectedUser, err := json.Marshal(tc.Response.User)
			if err != nil {
				t.Error(err)
			}

			jsonResponseUser, err := json.Marshal(getUserResponse.User)
			if err != nil {
				t.Error(err)
			}
			t.Fatalf("test case userId %d, expected response user %+v, got %+v", tc.UserId, string(jsonExpectedUser), string(jsonResponseUser))
		}
	}
}
