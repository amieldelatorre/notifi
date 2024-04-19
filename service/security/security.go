package security // import "github.com/amieldelatorre/notifi/service/security"

import (
	"context"
	"log/slog"

	"github.com/alexedwards/argon2id"
	"github.com/amieldelatorre/notifi/middleware"
)

func HashPassword(ctx context.Context, password string, logger *slog.Logger) (string, error) {
	requestId := ctx.Value(middleware.RequestIdName)
	logger.Debug("Creating user", "requestId", requestId)

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func IsCorrectPassword(ctx context.Context, password string, hash string, logger *slog.Logger) (bool, error) {
	requestId := ctx.Value(middleware.RequestIdName)
	logger.Debug("Creating user", "requestId", requestId)

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, err
}
