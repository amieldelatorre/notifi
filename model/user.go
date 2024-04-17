package model // import "github.com/amieldelatorre/notifi/model"

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
