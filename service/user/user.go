package user // import "github.com/amieldelatorre/notifi/service/user"

import (
	"context"
	"errors"
	"net/http"

	"github.com/amieldelatorre/notifi/model"
	userProvider "github.com/amieldelatorre/notifi/repository/user"
	"github.com/jackc/pgx/v5"
)

type GetUserResponse struct {
	Success bool              `json:"success"`
	User    *model.User       `json:"user,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type Service struct {
	Provider userProvider.UserProvider
}

func New(provider userProvider.UserProvider) Service {
	return Service{Provider: provider}
}

func (service *Service) GetUserById(ctx context.Context, id int) (int, GetUserResponse) {
	response := GetUserResponse{
		Success: false,
		Errors:  make(map[string]string),
	}

	user, err := service.Provider.GetUserById(ctx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		response.Errors["user"] = "User not found"
		return http.StatusNotFound, response
	} else if err != nil {
		response.Errors["server"] = "Something went wrong"
		return http.StatusNotFound, response
	}

	response.Success = true
	response.User = &user

	return http.StatusOK, response
}
