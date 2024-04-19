package security

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/amieldelatorre/notifi/logger"
)

// Would be difficult to test the actual hash of the password since it gives a different hash
// Just test the hashing and check of 1 password
func TestPasswordHashingAndChecking(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	tcs := []struct {
		ExpectedMatch    bool
		OriginalPassword string
		PasswordToCheck  string
	}{
		{
			ExpectedMatch:    true,
			OriginalPassword: "password123",
			PasswordToCheck:  "password123",
		},
		{
			ExpectedMatch:    false,
			OriginalPassword: "password123",
			PasswordToCheck:  "123password",
		},
	}

	for _, tc := range tcs {
		hashedPassword, err := HashPassword(context.Background(), tc.OriginalPassword, logger)
		if err != nil {
			t.Error(err)
		}

		passwordMatch, err := IsCorrectPassword(context.Background(), tc.PasswordToCheck, hashedPassword, logger)
		if err != nil {
			t.Error(err)
		}

		if passwordMatch != tc.ExpectedMatch {
			t.Fatalf("expected passwords match result '%t', got '%t'", tc.ExpectedMatch, passwordMatch)
		}
	}
}
