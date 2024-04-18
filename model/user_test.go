package model

import (
	"reflect"
	"testing"
)

type UserInputTestCases struct {
	ExpectedResponseErrors map[string][]string
	UserInput              UserInput
	ExpectedCleanInput     UserInput
}

func TestUserInputValidate(t *testing.T) {
	testCases := GetUserInputTestCases()

	for _, tc := range testCases {
		actualCleanInput, actualResponse := tc.UserInput.Validate()
		if !reflect.DeepEqual(tc.ExpectedResponseErrors, actualResponse) {
			t.Fatalf("expected %+v, got %+v", tc.ExpectedResponseErrors, actualResponse)
		}

		if !reflect.DeepEqual(tc.ExpectedCleanInput, actualCleanInput) {
			t.Fatalf("expected %+v, got %+v", tc.ExpectedCleanInput, actualCleanInput)
		}
	}
}

func GetUserInputTestCases() []UserInputTestCases {
	tc1 := UserInputTestCases{
		UserInput: UserInput{
			Email:     "",
			FirstName: "",
			LastName:  "",
			Password:  "",
		},
		ExpectedResponseErrors: map[string][]string{
			"email":     {"Cannot be empty"},
			"firstName": {"Cannot be empty"},
			"lastName":  {"Cannot be empty"},
			"password":  {"Cannot be empty and must be at least 8 characters"},
		},
		ExpectedCleanInput: UserInput{
			Email:     "",
			FirstName: "",
			LastName:  "",
			Password:  "",
		},
	}

	tc2 := UserInputTestCases{
		UserInput: UserInput{
			Email:     "isaac.newton@example.invalid",
			FirstName: "Isaac",
			LastName:  "Newton",
			Password:  "Password",
		},
		ExpectedResponseErrors: map[string][]string{},
		ExpectedCleanInput: UserInput{
			Email:     "isaac.newton@example.invalid",
			FirstName: "Isaac",
			LastName:  "Newton",
			Password:  "Password",
		},
	}

	tc3 := UserInputTestCases{
		UserInput: UserInput{
			Email:     "isaac.newton@",
			FirstName: "Isaac",
			LastName:  "Newton",
			Password:  "Password",
		},
		ExpectedResponseErrors: map[string][]string{
			"email": {"Must be a valid format: mail@example.invalid"},
		},
		ExpectedCleanInput: UserInput{
			Email:     "isaac.newton@",
			FirstName: "Isaac",
			LastName:  "Newton",
			Password:  "Password",
		},
	}

	tc4 := UserInputTestCases{
		UserInput: UserInput{
			Email:     "isaac.newton@invalid.com",
			FirstName: "Isaac      ",
			LastName:  "Newton",
			Password:  "Password",
		},
		ExpectedResponseErrors: map[string][]string{},
		ExpectedCleanInput: UserInput{
			Email:     "isaac.newton@invalid.com",
			FirstName: "Isaac",
			LastName:  "Newton",
			Password:  "Password",
		},
	}

	tcs := []UserInputTestCases{tc1, tc2, tc3, tc4}
	return tcs
}
