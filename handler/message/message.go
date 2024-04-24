package message

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/middleware"
	"github.com/amieldelatorre/notifi/model"
	messageService "github.com/amieldelatorre/notifi/service/message"
	"github.com/amieldelatorre/notifi/service/security"
	"github.com/amieldelatorre/notifi/utils"
)

type MessageHandler struct {
	Logger     *slog.Logger
	Service    messageService.Service
	JwtService security.JwtService
}

func New(logger *slog.Logger, service messageService.Service, jwtService security.JwtService) MessageHandler {
	return MessageHandler{Logger: logger, Service: service, JwtService: jwtService}
}

func (h *MessageHandler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("Registering routes for the message handler")
	m := middleware.New(h.Logger, h.JwtService)
	postMessageHandler := m.RecoverPanic(m.AddRequestId(m.RequireJwtToken(m.RequireApplicationJson(http.HandlerFunc(h.postMessage)))))

	mux.Handle("POST /api/v1/message", postMessageHandler)
}

func (h *MessageHandler) postMessage(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestIdName)
	h.Logger.Debug("Creating message", "requestId", requestId)
	w.Header().Set("Content-Type", "application/json")

	var messageInput model.MessageInput
	var response messageService.Response

	err := json.NewDecoder(r.Body).Decode(&messageInput)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); ok {
			h.Logger.Error("Post Message, could not unmarshal json from request body", "requestId", requestId, "error", err, "responseStatusCode", http.StatusInternalServerError)
			response.Errors["server"] = append(response.Errors["server"], "Something went wrong")

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			h.Logger.Info("Post Message", "requestId", requestId, "responseStatusCode", http.StatusBadRequest)
			response.Errors["messageInput"] = append(response.Errors["messageInput"], "Invalid json")

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	statusCode, response := h.Service.CreateMessage(r.Context(), messageInput)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Post Message", "requestId", requestId, "responseStatusCode", statusCode)
}
