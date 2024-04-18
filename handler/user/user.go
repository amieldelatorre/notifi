package user // import "github.com/amieldelatorre/notifi/handler/user"

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amieldelatorre/notifi/middleware"
	userService "github.com/amieldelatorre/notifi/service/user"
)

type UserHandler struct {
	Logger  *slog.Logger
	Service userService.Service
}

func New(logger *slog.Logger, service userService.Service) UserHandler {
	return UserHandler{Logger: logger, Service: service}
}

func (h UserHandler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("Registering routes for the user handler")
	m := middleware.New(h.Logger)
	getUserHandler := m.AddRequestId(http.HandlerFunc(h.getUser))

	mux.HandleFunc("POST /api/v1/user", h.postUser)
	mux.Handle("GET /api/v1/user", getUserHandler)
	mux.HandleFunc("PUT /api/v1/user", h.putUser)
	mux.HandleFunc("DELETE /api/v1/user", h.deleteUser)

	// mux.Handle("GET /api/v1/user", middleware.ApiKeyAuth(getUserHandler))
}

func (h UserHandler) postUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}

func (h UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := strconv.Atoi(r.Header.Get("x-user-id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	statusCode, response := h.Service.GetUserById(r.Context(), userId)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h UserHandler) putUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}

func (h UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}
