package user // import "github.com/amieldelatorre/notifi/backend/handler/user"

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/backend/middleware"
	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/amieldelatorre/notifi/backend/service/security"
	userService "github.com/amieldelatorre/notifi/backend/service/user"
	"github.com/amieldelatorre/notifi/backend/utils"
)

type UserHandler struct {
	Logger     *slog.Logger
	Service    userService.Service
	JwtService security.JwtService
}

func New(logger *slog.Logger, service userService.Service, jwtService security.JwtService) UserHandler {
	return UserHandler{Logger: logger, Service: service, JwtService: jwtService}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("Registering routes for the user handler")
	m := middleware.New(h.Logger, h.JwtService)
	getUserHandler := m.RecoverPanic(m.AddRequestId(m.RequireJwtToken(http.HandlerFunc(h.getUser))))
	postUserHandler := m.RecoverPanic(m.AddRequestId(m.RequireApplicationJson(http.HandlerFunc(h.postUser))))

	mux.Handle("POST /api/v1/user", postUserHandler)
	mux.Handle("GET /api/v1/user", getUserHandler)
	mux.HandleFunc("PUT /api/v1/user", h.putUser)
	mux.HandleFunc("DELETE /api/v1/user", h.deleteUser)
}

func (h *UserHandler) postUser(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestIdName)
	h.Logger.Debug("Creating user", "requestId", requestId)

	var userInput model.UserInput
	var response userService.UserResponse

	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); ok {
			h.Logger.Error("Post User, could not unmarshal json from request body", "requestId", requestId, "error", err, "responseStatusCode", http.StatusInternalServerError)

			response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
			utils.EncodeResponse[userService.UserResponse](w, http.StatusInternalServerError, response)
			return
		} else {
			h.Logger.Info("Post User", "requestId", requestId, "responseStatusCode", http.StatusBadRequest)

			response.Errors["userInput"] = append(response.Errors["userInput"], "Invalid json")
			utils.EncodeResponse[userService.UserResponse](w, http.StatusBadRequest, response)
			return
		}
	}

	statusCode, response := h.Service.CreateUser(r.Context(), userInput)
	utils.EncodeResponse[userService.UserResponse](w, statusCode, response)
	h.Logger.Info("Post User", "requestId", requestId, "responseStatusCode", statusCode)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestIdName)
	h.Logger.Debug("Retrieving user", "requestId", requestId)

	// TODO: Use context to get user id instead!
	userId := r.Context().Value(utils.UserId)

	statusCode, response := h.Service.GetUserById(r.Context(), userId.(int))
	utils.EncodeResponse[userService.UserResponse](w, statusCode, response)
	h.Logger.Info("Get User", "requestId", requestId, "responseStatusCode", statusCode)
}

func (h *UserHandler) putUser(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotImplemented)
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotImplemented)
}
