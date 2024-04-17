package user // import "github.com/amieldelatorre/notifi/handler/user"

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	userService "github.com/amieldelatorre/notifi/service/user"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	Service userService.Service
}

func New(service userService.Service) UserHandler {
	return UserHandler{Service: service}
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

	user, err := h.Service.GetUserById(context.Background(), userId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h UserHandler) putUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}

func (h UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}
