package model // import "github.com/amieldelatorre/notifi/model"

import "time"

type User struct {
	Id              int       `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Password        string    `json:"-"`
	DatetimeCreated time.Time `json:"datetimeCreated"`
	DatetimeUpdated time.Time `json:"datetimeUpdated"`
}
