package user // import "github.com/amieldelatorre/notifi/repository/user"

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
)

type UserProvider interface {
	GetUserById(ctx context.Context, id int) (model.User, error)
}

type UserDb struct {
	Provider UserProvider
}

func NewUserDb(provider UserProvider) *UserDb {
	return &UserDb{Provider: provider}
}
