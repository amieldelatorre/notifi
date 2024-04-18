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

type UserResponse struct {
	User   *model.User         `json:"user,omitempty"`
	Errors map[string][]string `json:"errors,omitempty"`
}

type Service struct {
	Provider userProvider.UserProvider
	Logger   *slog.Logger
}

func New(logger *slog.Logger, provider userProvider.UserProvider) Service {
	return Service{Logger: logger, Provider: provider}
}

func (service *Service) CreateUser(ctx context.Context, input model.UserInput) (int, UserResponse) {
	response := UserResponse{
		Errors: make(map[string][]string),
	}

	cleanInput, validationErrs, err := service.validateUserinput(ctx, input)
	if err != nil {
		service.Logger.Error("Could not check if email exists", "requestId", ctx.Value(middleware.RequestIdName), "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	} else if len(validationErrs) != 0 {
		response.Errors = validationErrs
		return http.StatusBadRequest, response
	}

	generatedId, err := service.Provider.CreateUser(ctx, cleanInput)
	if err != nil {
		service.Logger.Error("Could not create user from provider", "requestId", ctx.Value(middleware.RequestIdName), "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	newUser, err := service.Provider.GetUserById(ctx, generatedId)
	if err != nil {
		service.Logger.Error("New user created but something went wrong with retrieving the new user", "requestId", ctx.Value(middleware.RequestIdName), "error", err)
		response.Errors["server"] = append(response.Errors["server"], "New user created but something went wrong with retrieving the new user")
		return http.StatusInternalServerError, response
	}

	response.User = &newUser

	return http.StatusCreated, response
}

func (service *Service) GetUserById(ctx context.Context, id int) (int, UserResponse) {
	response := UserResponse{
		Errors: make(map[string][]string),
	}

	user, err := service.Provider.GetUserById(ctx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		response.Errors["user"] = append(response.Errors["user"], "User not found")
		return http.StatusNotFound, response
	} else if err != nil {
		service.Logger.Error("Could not get user from provider", "requestId", ctx.Value(middleware.RequestIdName), "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	response.User = &user

	return http.StatusOK, response
}

func (service *Service) EmailExists(ctx context.Context, email string) (bool, error) {
	_, err := service.Provider.GetUserByEmail(ctx, email)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (service *Service) validateUserinput(ctx context.Context, input model.UserInput) (model.UserInput, map[string][]string, error) {
	cleanInput, validationErrs := input.Validate()

	emailExists, err := service.EmailExists(ctx, cleanInput.Email)
	if err != nil {
		return cleanInput, validationErrs, err
	}
	if emailExists {
		validationErrs["email"] = append(validationErrs["email"], "Email already exists")
	}

	return cleanInput, validationErrs, nil
}
