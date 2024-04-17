package user

import (
	"context"
	"net/http"
	"time"

	"github.com/amieldelatorre/notifi/model"
	"github.com/jackc/pgx/v5"
)

type MockUserProvider struct {
	Users []model.User
}

type TestGetUserByIdTestCase struct {
	ExpectedStatusCode int
	Response           GetUserResponse
	UserId             int
}

func GetTestUsers() []model.User {
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
	return MockUserProvider{Users: GetTestUsers()}
}

func (mr *MockUserProvider) GetUserById(ctx context.Context, id int) (model.User, error) {
	for _, val := range mr.Users {
		if val.Id == id {
			return val, nil
		}
	}

	return model.User{}, pgx.ErrNoRows
}

func GetValidTestGetUserByIdTestCases() []TestGetUserByIdTestCase {
	testCases := []TestGetUserByIdTestCase{}

	for _, u := range GetTestUsers() {
		tc := TestGetUserByIdTestCase{
			ExpectedStatusCode: http.StatusOK,
			Response:           GetUserResponse{Success: true, User: &u, Errors: map[string]string{}},
			UserId:             u.Id,
		}
		testCases = append(testCases, tc)
	}
	return testCases
}

func GetInvalidTestGetUserByIdTestCase() TestGetUserByIdTestCase {
	return TestGetUserByIdTestCase{
		ExpectedStatusCode: http.StatusNotFound,
		Response:           GetUserResponse{Success: false, Errors: map[string]string{"user": "User not found"}},
		UserId:             100,
	}
}
