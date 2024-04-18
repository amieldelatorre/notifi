package user // import "github.com/amieldelatorre/notifi/handler/user"

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amieldelatorre/notifi/middleware"
	"github.com/amieldelatorre/notifi/model"
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
	postUserHandler := m.AddRequestId(http.HandlerFunc(h.postUser))

	mux.Handle("POST /api/v1/user", postUserHandler)
	mux.Handle("GET /api/v1/user", getUserHandler)
	mux.HandleFunc("PUT /api/v1/user", h.putUser)
	mux.HandleFunc("DELETE /api/v1/user", h.deleteUser)

	// mux.Handle("GET /api/v1/user", middleware.ApiKeyAuth(getUserHandler))
}

func (h UserHandler) postUser(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(middleware.RequestIdName)
	h.Logger.Debug("Creating user", "requestId", requestId)
	w.Header().Set("Content-Type", "application/json")

	var userInput model.UserInput

	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); ok {
			h.Logger.Error("When creating user, could not unmarshal json", "requestId", requestId, "error", err, "responseStatusCode", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.Logger.Info("Post User", "requestId", requestId, "responseStatusCode", http.StatusBadRequest)
			return
		}
	}

	statusCode, response := h.Service.CreateUser(r.Context(), userInput)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Post User", "requestId", requestId, "responseStatusCode", statusCode)
}

func (h UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(middleware.RequestIdName)
	h.Logger.Debug("Retrieving user", "requestId", requestId)
	w.Header().Set("Content-Type", "application/json")

	// TODO: Use context to get user id instead!
	userId, err := strconv.Atoi(r.Header.Get("x-user-id"))
	if err != nil {
		h.Logger.Error("When retrieving user, could not convert string to int", "requestId", requestId, "error", err, "responseStatusCode", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	statusCode, response := h.Service.GetUserById(r.Context(), userId)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Get User", "requestId", requestId, "responseStatusCode", statusCode)
}

func (h UserHandler) putUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}

func (h UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
}
