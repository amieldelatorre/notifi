package user // import "github.com/amieldelatorre/notifi/repository/user"

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
)

type UserProvider interface {
	CreateUser(ctx context.Context, input model.UserInput) (int, error)
	GetUserById(ctx context.Context, id int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}
