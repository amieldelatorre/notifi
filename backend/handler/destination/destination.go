package destination // import "github.com/amieldelatorre/notifi/backend/handler/destination"

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amieldelatorre/notifi/backend/middleware"
	"github.com/amieldelatorre/notifi/backend/model"
	destinationService "github.com/amieldelatorre/notifi/backend/service/destination"
	"github.com/amieldelatorre/notifi/backend/service/security"
	"github.com/amieldelatorre/notifi/backend/utils"
)

type DestinationHandler struct {
	Logger     *slog.Logger
	Service    destinationService.Service
	JwtService security.JwtService
}

func New(logger *slog.Logger, service destinationService.Service, jwtService security.JwtService) DestinationHandler {
	return DestinationHandler{Logger: logger, Service: service, JwtService: jwtService}
}

func (h *DestinationHandler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("Registering routes for the destination handler")
	m := middleware.New(h.Logger, h.JwtService)
	postDestinationHandlerFunc := m.RecoverPanic(m.AddRequestId(m.RequireJwtToken(m.RequireApplicationJson(http.HandlerFunc(h.postDestination)))))
	getDestinationsHandlerFunc := m.RecoverPanic(m.AddRequestId(m.RequireJwtToken((http.HandlerFunc(h.getDestinations)))))
	getDestinationByIdHandlerFunc := m.RecoverPanic(m.AddRequestId(m.RequireJwtToken((http.HandlerFunc(h.getDestinationById)))))

	mux.Handle("POST /api/v1/destination", postDestinationHandlerFunc)
	mux.Handle("GET /api/v1/destination", getDestinationsHandlerFunc)
	mux.Handle("GET /api/v1/destination/{id}", getDestinationByIdHandlerFunc)
}

func (h *DestinationHandler) postDestination(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestIdName)
	h.Logger.Debug("Creating destination", "requestId", requestId)
	w.Header().Set("Content-Type", "application/json")

	var destinationInput model.DestinationInput
	var response destinationService.Response

	err := json.NewDecoder(r.Body).Decode(&destinationInput)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); ok {
			h.Logger.Error("Post Destination, could not unmarshal json from request body", "requestId", requestId, "error", err, "responseStatusCode", http.StatusInternalServerError)
			response.Errors["server"] = append(response.Errors["server"], "Something went wrong")

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			h.Logger.Info("Post Destination", "requestId", requestId, "responseStatusCode", http.StatusBadRequest)
			response.Errors["destinationInput"] = append(response.Errors["destinationInput"], "Invalid json")

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	statusCode, response := h.Service.CreateDestination(r.Context(), destinationInput)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Post Destination", "requestId", requestId, "responseStatusCode", statusCode)
}

func (h *DestinationHandler) getDestinations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(utils.RequestIdName)
	h.Logger.Debug("Retrieving destinations", "requestId", requestId)
	w.Header().Set("Content-Type", "application/json")

	userId := ctx.Value(utils.UserId).(int)

	statusCode, response := h.Service.GetAllDestinations(ctx, userId)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Get Destinations", "requestId", requestId, "responseStatusCode", statusCode)
}

func (h *DestinationHandler) getDestinationById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(utils.RequestIdName)
	h.Logger.Debug("Retrieving destinations", "requestId", requestId)
	w.Header().Set("Content-Type", "application/json")

	userId := ctx.Value(utils.UserId).(int)
	destinationIdString := r.PathValue("id")
	response := destinationService.Response{
		Errors: map[string][]string{},
	}

	destinationIdInt, err := strconv.Atoi(destinationIdString)
	if err != nil {
		h.Logger.Debug("Get Destination By Id", "requestId", requestId, "responseStatusCode", http.StatusBadRequest)
		response.Errors["id"] = append(response.Errors["id"], "Invalid destination Id provided. Must be an integer")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	statusCode, response := h.Service.GetDestinationById(ctx, destinationIdInt, userId)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	h.Logger.Info("Get Destination By Id", "requestId", requestId, "responseStatusCode", statusCode)
}
