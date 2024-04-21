package security

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAccessTokenSuccess(t *testing.T) {
	initialClaims := UserClaims{
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
	createdToken, err := CreateAccessToken(initialClaims, signingKey)
	if err != nil {
		t.Fatalf("Expected no error, got %+v", err)
	}

	retrievedClaims, err := ParseAccessToken(createdToken, signingKey)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(initialClaims, *retrievedClaims) {
		t.Fatalf("expected claims '%+v', got '%+v'", initialClaims, *retrievedClaims)
	}
}

func TestAccessTokenFail(t *testing.T) {
	tc := struct {
		ExpectedError  error
		ExpectedClaims *UserClaims
		token          string
	}{
		ExpectedError:  &InvalidAccessToken{},
		ExpectedClaims: nil,
		token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiJpc2FhYy5uZXd0b25AZXhhbXBsZS5pbnZhbGlkIiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiYXVkIjpbInNvbWVib2R5X2Vsc2UiXSwiZXhwIjozMjUwMzY4MDAwMCwibmJmIjozMjUwMzY4MDAwMCwiaWF0IjozMjUwMzY4MDAwMCwianRpIjoiMSJ9.uCsh2Yp2SdXeDtG_Vgdhv9rMuwJt4cj9Yx2cfP6P_lY",
	}

	signingKey := []byte("This_is_a_super_secret_key")
	claims, err := ParseAccessToken(tc.token, signingKey)
	if err != tc.ExpectedError {
		t.Fatalf("expected error '%s', got '%s'", tc.ExpectedError, err)
	}

	if claims != nil {
		t.Fatalf("expected claims as nil, got '%+v'", claims)
	}
}
