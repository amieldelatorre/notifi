package model // import "github.com/amieldelatorre/notifi/model"

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

const (
	passwordMinLength = 8
)

type User struct {
	Id              int       `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Password        string    `json:"-"`
	DatetimeCreated time.Time `json:"datetimeCreated"`
	DatetimeUpdated time.Time `json:"datetimeUpdated"`
}

type UserInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func (input *UserInput) Validate() (UserInput, map[string][]string) {
	validationErrors := make(map[string][]string)

	if strings.TrimSpace(input.FirstName) == "" {
		validationErrors["firstName"] = append(validationErrors["firstName"], "Cannot be empty")
	}

	if strings.TrimSpace(input.LastName) == "" {
		validationErrors["lastName"] = append(validationErrors["lastName"], "Cannot be empty")
	}

	if strings.TrimSpace(input.Email) == "" {
		validationErrors["email"] = append(validationErrors["email"], "Cannot be empty")
	} else if _, err := mail.ParseAddress(strings.TrimSpace(input.Email)); err != nil {
		validationErrors["email"] = append(validationErrors["email"], "Must be a valid format: mail@example.invalid")
	}

	if strings.TrimSpace(input.Password) == "" || len(strings.TrimSpace(input.Password)) < passwordMinLength {
		validationErrors["password"] = append(validationErrors["password"], fmt.Sprintf("Cannot be empty and must be at least %d characters", passwordMinLength))
	}

	cleanInput := UserInput{
		Email:     strings.TrimSpace(input.Email),
		FirstName: strings.TrimSpace(input.FirstName),
		LastName:  strings.TrimSpace(input.LastName),
		// Not trimming password because as long as there are 8 non white space characters that are in the middle of the string it satisfies requirements
		Password: input.Password,
	}

	return cleanInput, validationErrors
}
