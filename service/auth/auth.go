package auth // import "github.com/amieldelatorre/notifi/service/auth"

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/service/security"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthProvider interface {
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type AuthResponse struct {
	Token  string              `json:"token,omitempty"`
	Errors map[string][]string `json:"errors,omitempty"`
}

type Service struct {
	Logger     *slog.Logger
	Provider   AuthProvider
	JwtService security.JwtService
}

type BasicAuthCredentials struct {
	Email    string
	Password string
}

func New(logger *slog.Logger, provider AuthProvider, jwtService security.JwtService) Service {
	return Service{Logger: logger, Provider: provider, JwtService: jwtService}
}

func (s *Service) LoginUser(ctx context.Context, basicAuthCredentials BasicAuthCredentials) (int, AuthResponse) {
	response := AuthResponse{
		Errors: make(map[string][]string),
	}
	requestId := ctx.Value(utils.RequestIdName)

	user, err := s.Provider.GetUserByEmail(ctx, basicAuthCredentials.Email)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		response.Errors["credentials"] = append(response.Errors["credentials"], "email and password combination not found.")
		return http.StatusForbidden, response
	} else if err != nil {
		s.Logger.Error("Could not get user from provider", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	validLogin, err := security.IsCorrectPassword(ctx, basicAuthCredentials.Password, user.Password, s.Logger)
	if err != nil {
		s.Logger.Error("Could not check if login is valid", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	if !validLogin {
		response.Errors["credentials"] = append(response.Errors["credentials"], "email/password combination not found.")
		return http.StatusForbidden, response
	}

	timeNow := time.Now()
	claimsUuid, err := uuid.NewV7()
	if err != nil {
		s.Logger.Error("Problem generating uuid for claims: %s", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	claims := security.UserClaims{
		user.Id,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
			NotBefore: jwt.NewNumericDate(timeNow),
			Issuer:    "Notifi",
			Subject:   "AccessToken",
			ID:        claimsUuid.String(),
			Audience:  []string{"Notifi"},
		},
	}

	token, err := s.JwtService.CreateAccessToken(claims)
	if err != nil {
		s.Logger.Error("Problem generating token for claims: %s", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	response.Token = token
	return http.StatusOK, response
}
