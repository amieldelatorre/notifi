package user // import "github.com/amieldelatorre/notifi/service/user"

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
	userProvider "github.com/amieldelatorre/notifi/repository/user"
)

type Service struct {
	Provider userProvider.UserProvider
}

func New(provider userProvider.UserProvider) Service {
	return Service{Provider: provider}
}

func (service *Service) GetUserById(ctx context.Context, id int) (model.User, error) {
	user, err := service.Provider.GetUserById(ctx, id)
	return user, err
}
