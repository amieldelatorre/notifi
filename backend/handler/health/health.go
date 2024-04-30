package health

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/backend/middleware"
	"github.com/amieldelatorre/notifi/backend/service/security"
	"github.com/amieldelatorre/notifi/backend/utils"
)

type DbProvider interface {
	IsHealthy(ctx context.Context) bool
}

type QueueProvider interface {
	IsHealthy(ctx context.Context) bool
}

type HealthHandler struct {
	Logger        *slog.Logger
	DbProvider    DbProvider
	QueueProvider QueueProvider
	JwtService    security.JwtService
}

type Response struct {
	Errors map[string][]string `json:"errors,omitempty"`
}

func New(logger *slog.Logger, dbProvider DbProvider, queueProvider QueueProvider, jwtService security.JwtService) HealthHandler {
	return HealthHandler{Logger: logger, DbProvider: dbProvider, QueueProvider: queueProvider, JwtService: jwtService}
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("Registering routes for the health handler")
	m := middleware.New(h.Logger, h.JwtService)
	healthCheckHandler := m.RecoverPanic(m.AddRequestId(http.HandlerFunc(h.healthCheck)))

	mux.Handle("GET /api/v1/health", healthCheckHandler)
}

func (h *HealthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestIdName)
	h.Logger.Debug("Health check", "requestId", requestId)

	statusCode := http.StatusOK
	response := Response{
		Errors: map[string][]string{},
	}

	if !h.DbProvider.IsHealthy(r.Context()) {
		response.Errors["database"] = append(response.Errors["database"], "Database is unhealthy")
	}

	if !h.QueueProvider.IsHealthy(r.Context()) {
		response.Errors["queue"] = append(response.Errors["queue"], "Message Queue is unhealthy")
	}

	if len(response.Errors) != 0 {
		statusCode = http.StatusInternalServerError
	}

	utils.EncodeResponse[Response](w, statusCode, response)
}
