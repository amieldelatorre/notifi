package auth // import "github.com/amieldelatorre/notifi/repository/auth"

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
)

type AuthProvider interface {
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}
