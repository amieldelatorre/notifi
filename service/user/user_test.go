package user

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/model"
	"github.com/jackc/pgx/v5"
)

type MockUserProvider struct {
	Users []model.User
}

func getTestUsers() []model.User {
	user1 := model.User{
		Id:              1,
		Email:           "isaac.newton@example.invalid",
		FirstName:       "Isaac",
		LastName:        "Newton",
		Password:        "Password",
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	}

	user2 := model.User{
		Id:              2,
		Email:           "alberteinstein@example.invalid",
		FirstName:       "Albert",
		LastName:        "Einstein",
		Password:        "Password",
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	}
	users := []model.User{user1, user2}
	return users
}

func NewMockUserRepo() MockUserProvider {
	return MockUserProvider{Users: getTestUsers()}
}

func (mr *MockUserProvider) GetUserById(ctx context.Context, id int) (model.User, error) {
	for _, val := range mr.Users {
		if val.Id == id {
			return val, nil
		}
	}

	return model.User{}, pgx.ErrNoRows
}

type TestGetUserByIdTestCase struct {
	ExpectedStatusCode int
	Response           GetUserResponse
	UserId             int
}

func TestGetUserById(t *testing.T) {
	mockUserProvider := NewMockUserRepo()
	service := New(&mockUserProvider)

	testCases := []TestGetUserByIdTestCase{}

	for _, u := range getTestUsers() {
		tc := TestGetUserByIdTestCase{
			ExpectedStatusCode: http.StatusOK,
			Response:           GetUserResponse{Success: true, User: &u, Errors: map[string]string{}},
			UserId:             u.Id,
		}
		testCases = append(testCases, tc)
	}

	testCases = append(testCases, TestGetUserByIdTestCase{
		ExpectedStatusCode: http.StatusNotFound,
		Response:           GetUserResponse{Success: false, Errors: map[string]string{"user": "User not found"}},
		UserId:             100,
	})

	for _, tc := range testCases {
		actualStatusCode, actualResponse := service.GetUserById(context.Background(), tc.UserId)
		if actualResponse.Success != tc.Response.Success {
			t.Fatalf("test case userId %d, expected response success %t, got %t", tc.UserId, tc.Response.Success, actualResponse.Success)
		}

		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case userId %d, expected exected status code %d, got %d", tc.UserId, tc.ExpectedStatusCode, actualStatusCode)
		}

		if len(actualResponse.Errors) != len(tc.Response.Errors) {
			t.Fatalf("expected response errors length %d, got %d", len(tc.Response.Errors), len(actualResponse.Errors))
		} else if !reflect.DeepEqual(actualResponse.Errors, tc.Response.Errors) {
			t.Fatalf("test case userId %d, expected response errors %+v, got %+v", tc.UserId, tc.Response.Errors, actualResponse.Errors)
		}

		if !reflect.DeepEqual(actualResponse.User, tc.Response.User) {
			t.Fatalf("test case userId %d, expected response user %+v, got %+v", tc.UserId, tc.Response.User, actualResponse.User)
		}
	}
}
