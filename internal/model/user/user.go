package user // import "github.com/amieldelatorre/notifi/internal/model/user"

import "time"

type User struct {
	Id              int
	Email           string
	FirstName       string
	LastName        string
	Password        string
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
}
