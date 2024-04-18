package user // import "github.com/amieldelatorre/notifi/service/user"

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/middleware"
	"github.com/amieldelatorre/notifi/model"
	userProvider "github.com/amieldelatorre/notifi/repository/user"
	"github.com/jackc/pgx/v5"
)

type GetUserResponse struct {
	User   *model.User       `json:"user,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

type Service struct {
	Provider userProvider.UserProvider
	Logger   *slog.Logger
}

func New(logger *slog.Logger, provider userProvider.UserProvider) Service {
	return Service{Logger: logger, Provider: provider}
}

func (service *Service) GetUserById(ctx context.Context, id int) (int, GetUserResponse) {
	response := GetUserResponse{
		Errors: make(map[string]string),
	}

	user, err := service.Provider.GetUserById(ctx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		response.Errors["user"] = "User not found"
		return http.StatusNotFound, response
	} else if err != nil {
		service.Logger.Error("Could not get user from provider", "requestId", ctx.Value(middleware.RequestIdName), "error", err)
		response.Errors["server"] = "Something went wrong"
		return http.StatusInternalServerError, response
	}

	response.User = &user

	return http.StatusOK, response
}
