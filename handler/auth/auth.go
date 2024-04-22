package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/middleware"
	authService "github.com/amieldelatorre/notifi/service/auth"
	"github.com/amieldelatorre/notifi/utils"
)

type AuthHandler struct {
	Logger  *slog.Logger
	Service authService.Service
}

func New(logger *slog.Logger, service authService.Service) AuthHandler {
	return AuthHandler{Logger: logger, Service: service}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("Registering routes for the auth handler")
	m := middleware.New(h.Logger)

	loginHandler := m.AddRequestId(http.HandlerFunc(h.login))

	mux.Handle("POST /api/v1/auth/login", loginHandler)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestIdName)
	h.Logger.Debug("Login user", "requestId", requestId)

	var basicAuthCredentials authService.BasicAuthCredentials

	err := json.NewDecoder(r.Body).Decode(&basicAuthCredentials)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); ok {
			h.Logger.Error("Login User, could not unmarshal json from request body", "requestId", requestId, "error", err, "responseStatusCode", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Info("Login User", "requestId", requestId, "responseStatusCode", http.StatusBadRequest)
			return
		}
	}

	statusCode, response := h.Service.LoginUser(r.Context(), basicAuthCredentials)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Login User", "requestId", requestId, "responseStatusCode", statusCode)
}
