package security

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAccessToken(t *testing.T) {
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiJpc2FhYy5uZXd0b25AZXhhbXBsZS5pbnZhbGlkIiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiYXVkIjpbInNvbWVib2R5X2Vsc2UiXSwiZXhwIjozMjUwMzYzMzIwMCwibmJmIjotMzA2MTAyNjcyMDAsImlhdCI6OTQ2NjM4MDAwLCJqdGkiOiIxIn0.pDv1N9-ji1WOq9Zx5EgtDS-2o_3I9viU0ntVRzMtQew"
	claims := UserClaims{
		1,
		"isaac.newton@example.invalid",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Date(3000, time.January, 1, 00, 0, 0, 0, time.Local)),
			IssuedAt:  jwt.NewNumericDate(time.Date(2000, time.January, 1, 00, 0, 0, 0, time.Local)),
			NotBefore: jwt.NewNumericDate(time.Date(1000, time.January, 1, 00, 0, 0, 0, time.Local)),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	signingKey := []byte("This_is_a_super_secret_key")
	token, err := CreateAccessToken(claims, signingKey)
	if err != nil {
		t.Error(err)
	}

	if token != expectedToken {
		t.Fatalf("expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestParseAccessToken(t *testing.T) {
	tcs := []struct {
		ExpectedError  error
		ExpectedClaims *UserClaims
		token          string
	}{
		{
			ExpectedError: nil,
			ExpectedClaims: &UserClaims{
				1,
				"isaac.newton@example.invalid",
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Date(3000, time.January, 1, 00, 0, 0, 0, time.Local)),
					IssuedAt:  jwt.NewNumericDate(time.Date(2000, time.January, 1, 00, 0, 0, 0, time.Local)),
					NotBefore: jwt.NewNumericDate(time.Date(1000, time.January, 1, 00, 0, 0, 0, time.Local)),
					Issuer:    "test",
					Subject:   "somebody",
					ID:        "1",
					Audience:  []string{"somebody_else"},
				},
			},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiJpc2FhYy5uZXd0b25AZXhhbXBsZS5pbnZhbGlkIiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiYXVkIjpbInNvbWVib2R5X2Vsc2UiXSwiZXhwIjozMjUwMzYzMzIwMCwibmJmIjotMzA2MTAyNjcyMDAsImlhdCI6OTQ2NjM4MDAwLCJqdGkiOiIxIn0.pDv1N9-ji1WOq9Zx5EgtDS-2o_3I9viU0ntVRzMtQew",
		},
		{
			ExpectedError:  &InvalidAccessToken{},
			ExpectedClaims: nil,
			token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiJpc2FhYy5uZXd0b25AZXhhbXBsZS5pbnZhbGlkIiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiYXVkIjpbInNvbWVib2R5X2Vsc2UiXSwiZXhwIjozMjUwMzY4MDAwMCwibmJmIjozMjUwMzY4MDAwMCwiaWF0IjozMjUwMzY4MDAwMCwianRpIjoiMSJ9.uCsh2Yp2SdXeDtG_Vgdhv9rMuwJt4cj9Yx2cfP6P_lY",
		},
	}

	signingKey := []byte("This_is_a_super_secret_key")

	for _, tc := range tcs {
		claims, err := ParseAccessToken(tc.token, signingKey)
		if err != tc.ExpectedError {
			t.Error(err)
		}

		if tc.ExpectedError == nil && claims == nil {
			t.Fatalf("test case '%s', claims is nil even though expected error is nil", tc.token)
		}

		// if tc.ExpectedClaims != claims {
		// 	t.Fatalf("test case '%s', expected claims '%+v', got '%+v'", tc.token, tc.ExpectedClaims, *claims)
		// }
		if !reflect.DeepEqual(tc.ExpectedClaims, claims) {
			t.Fatalf("test case '%s', expected claims '%+v', got '%+v'", tc.token, tc.ExpectedClaims, *claims)
		}

	}

}
