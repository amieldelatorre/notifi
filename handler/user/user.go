package user // import "github.com/amieldelatorre/notifi/handler/user"

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/amieldelatorre/notifi/logger"
	userService "github.com/amieldelatorre/notifi/service/user"
)

type UserHandler struct {
	Logger  *logger.Logger
	Service userService.Service
}

func New(logger *logger.Logger, service userService.Service) UserHandler {
	return UserHandler{Logger: logger, Service: service}
}

func (h UserHandler) RegisterRoutes(mux *http.ServeMux) {
	getUserHandler := http.HandlerFunc(h.getUser)

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
