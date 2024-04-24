package model

import (
	"reflect"
	"testing"
)

type DestinationInputTestCase struct {
	ExpectedResponseErrors map[string][]string
	DestinationInput       DestinationInput
	ExpectedCleanInput     DestinationInput
}

func TestDestinationInputValidate(t *testing.T) {
	testCases := GetDestinationInputTestCases()

	for _, tc := range testCases {
		actualCleanInput, actualResponse := tc.DestinationInput.Validate()
		if !reflect.DeepEqual(tc.ExpectedResponseErrors, actualResponse) {
			t.Fatalf("expected %+v, got %+v", tc.ExpectedResponseErrors, actualResponse)
		}

		if !reflect.DeepEqual(tc.ExpectedCleanInput, actualCleanInput) {
			t.Fatalf("expected %+v, got %+v", tc.ExpectedCleanInput, actualCleanInput)
		}
	}
}

func GetDestinationInputTestCases() []DestinationInputTestCase {
	tcs := []DestinationInputTestCase{
		{
			DestinationInput: DestinationInput{
				Type:       " ",
				Identifier: "",
			},
			ExpectedResponseErrors: map[string][]string{
				"type":       {"Must be one of DISCORD"},
				"identifier": {"Cannot be empty"},
			},
			ExpectedCleanInput: DestinationInput{
				Type:       "",
				Identifier: "",
			},
		},
		{
			DestinationInput: DestinationInput{
				Type:       "x",
				Identifier: "anidentifier",
			},
			ExpectedResponseErrors: map[string][]string{
				"type": {"Must be one of DISCORD"},
			},
			ExpectedCleanInput: DestinationInput{
				Type:       "X",
				Identifier: "anidentifier",
			},
		},
		{
			DestinationInput: DestinationInput{
				Type:       "DISCORD ",
				Identifier: "anidentifier",
			},
			ExpectedResponseErrors: map[string][]string{},
			ExpectedCleanInput: DestinationInput{
				Type:       "DISCORD",
				Identifier: "anidentifier",
			},
		},
		{
			DestinationInput: DestinationInput{
				Type:       "discord",
				Identifier: "anidentifier",
			},
			ExpectedResponseErrors: map[string][]string{},
			ExpectedCleanInput: DestinationInput{
				Type:       "DISCORD",
				Identifier: "anidentifier",
			},
		},
		{
			DestinationInput: DestinationInput{
				Type:       "DISCORD",
				Identifier: "anidentifier",
			},
			ExpectedResponseErrors: map[string][]string{},
			ExpectedCleanInput: DestinationInput{
				Type:       "DISCORD",
				Identifier: "anidentifier",
			},
		},
	}

	return tcs
}
