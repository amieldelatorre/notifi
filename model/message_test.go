package model

import (
	"reflect"
	"testing"
)

type MessageInputTestCase struct {
	ExpectedResponseErrors map[string][]string
	MessageInput           MessageInput
	ExpectedCleanInput     MessageInput
}

func TestMessageInputValidate(t *testing.T) {
	testCases := GetMessageInputTestCases()

	for _, tc := range testCases {
		actualCleanInput, actualResponse := tc.MessageInput.Validate()
		if !reflect.DeepEqual(tc.ExpectedResponseErrors, actualResponse) {
			t.Fatalf("expected %+v, got %+v", tc.ExpectedResponseErrors, actualResponse)
		}

		if !reflect.DeepEqual(tc.ExpectedCleanInput, actualCleanInput) {
			t.Fatalf("expected %+v, got %+v", tc.ExpectedCleanInput, actualCleanInput)
		}
	}
}

func GetMessageInputTestCases() []MessageInputTestCase {
	tc1 := MessageInputTestCase{
		MessageInput: MessageInput{
			Title: "",
			Body:  "",
		},
		ExpectedResponseErrors: map[string][]string{
			"title": {"Must have at least one non-whitespace character"},
			"body":  {"Must have at least one non-whitespace character"},
		},
		ExpectedCleanInput: MessageInput{
			Title: "",
			Body:  "",
		},
	}

	tc2 := MessageInputTestCase{
		MessageInput: MessageInput{
			Title: "s",
			Body:  "s",
		},
		ExpectedResponseErrors: map[string][]string{},
		ExpectedCleanInput: MessageInput{
			Title: "s",
			Body:  "s",
		},
	}

	tcs := []MessageInputTestCase{tc1, tc2}
	return tcs
}
