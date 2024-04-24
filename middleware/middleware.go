package middleware // import "github.com/amieldelatorre/notifi/middleware"

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/amieldelatorre/notifi/service/security"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/google/uuid"
)

type MiddlewareErrors struct {
	Errors map[string][]string `json:"errors,omitempty"`
}

type Middleware struct {
	Logger     *slog.Logger
	JwtService *security.JwtService
}

func New(logger *slog.Logger, jwtService security.JwtService) Middleware {
	return Middleware{Logger: logger, JwtService: &jwtService}
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.Error("Had to recover from panic", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Logger.Debug("Adding a request id to incoming request")
		id, err := uuid.NewV7()
		if err != nil {
			m.Logger.Error("Problem generating uuid: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.Logger.Debug("Request id generated and added to context", "requestId", id.String())

		ctx := context.WithValue(r.Context(), utils.RequestIdName, id.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) RequireJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := MiddlewareErrors{
			Errors: make(map[string][]string),
		}
		requestId := r.Context().Value(utils.RequestIdName)
		m.Logger.Debug("Checking if request has the Authorization Bearer token", "requestId", requestId)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Errors["Authorization"] = append(response.Errors["Authorization"], "Missing Authorization header")
			m.Logger.Debug("Missing Authorization header for JWT", "requestId", requestId)

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		authHeaderSlice := strings.Split(authHeader, " ")
		if len(authHeaderSlice) != 2 {
			response.Errors["Authorization"] = append(response.Errors["Authorization"], "Wrong Authorization header type, this endpoint requires an Authorization of 'Bearer' and a token")
			m.Logger.Debug("Authorization header value does not contain the right amount of values (2)", "requestId", requestId)

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		authType := authHeaderSlice[0]
		token := authHeaderSlice[1]

		if authType != "Bearer" {
			response.Errors["Authorization"] = append(response.Errors["Authorization"], "Wrong Authorization header type, this endpoint requires an Authorization of 'Bearer'")
			m.Logger.Debug("Header type not 'Bearer' for JWT", "requestId", requestId)

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		userClaims, err := m.JwtService.ParseAccessToken(token)
		if err != nil {
			response.Errors["Authorization"] = append(response.Errors["Authorization"], "Invalid 'Bearer' token")
			m.Logger.Debug("Invalid JWT token", "requestId", requestId)

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserId, userClaims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) RequireApplicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := MiddlewareErrors{
			Errors: make(map[string][]string),
		}
		requestId := r.Context().Value(utils.RequestIdName)
		m.Logger.Debug("Checking if request Content-Type is 'application/json'", "requestId", requestId)

		contentTypeHeader := r.Header.Get("content-type")
		if contentTypeHeader != "application/json" {
			response.Errors["Content-Type"] = append(response.Errors["Content-Type"], "Content-Type unsupported, must be 'application/json'")
			m.Logger.Debug("Content-Type header is not 'application/json'", "requestId", requestId)

			w.WriteHeader(http.StatusUnsupportedMediaType)
			json.NewEncoder(w).Encode(response)
		}
	})
}
