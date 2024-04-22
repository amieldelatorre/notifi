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
	Response           UserResponse
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

func (mr *MockUserProvider) CreateUser(ctx context.Context, input model.UserInput) (model.User, error) {
	id := len(mr.Users) + 1
	newUser := model.User{
		Id:              id,
		Email:           input.Email,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Password:        input.Password,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	}

	mr.Users = append(mr.Users, newUser)
	return newUser, nil
}

func (mr *MockUserProvider) GetUserById(ctx context.Context, id int) (model.User, error) {
	for _, val := range mr.Users {
		if val.Id == id {
			return val, nil
		}
	}

	return model.User{}, pgx.ErrNoRows
}

func (mr *MockUserProvider) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	for _, user := range mr.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.User{}, pgx.ErrNoRows
}

func GetValidTestGetUserByIdTestCases() []TestGetUserByIdTestCase {
	testCases := []TestGetUserByIdTestCase{}

	for _, u := range GetTestUsers() {
		tc := TestGetUserByIdTestCase{
			ExpectedStatusCode: http.StatusOK,
			Response:           UserResponse{User: &u, Errors: map[string][]string{}},
			UserId:             u.Id,
		}
		testCases = append(testCases, tc)
	}
	return testCases
}

func GetInvalidTestGetUserByIdTestCase() TestGetUserByIdTestCase {
	return TestGetUserByIdTestCase{
		ExpectedStatusCode: http.StatusNotFound,
		Response:           UserResponse{Errors: map[string][]string{"user": {"User not found"}}},
		UserId:             100,
	}
}
