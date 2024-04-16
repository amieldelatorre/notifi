package user

import (
	"testing"
	"time"

	userModel "github.com/amieldelatorre/notifi/internal/model/user"
)

func TestGetUserById(t *testing.T) {

}

type TestPostgrestDb struct {
	Users []userModel.User
}

func (testDb TestPostgrestDb) GetUserById(userId int) (userModel.User, error) {
	return testDb.Users[0], nil
}

func InitServiceTestHelper() Service {
	var users []userModel.User
	users = append(users, userModel.User{
		Id:              1,
		Email:           "test1@example.invalid",
		FirstName:       "James",
		LastName:        "Smith",
		Password:        "password",
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	})

	testDb := TestPostgrestDb{Users: users}
	service := New(testDb)
	return service
}
