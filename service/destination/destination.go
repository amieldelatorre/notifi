package destination // import "github.com/amieldelatorre/notifi/service/destination"

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/utils"
)

type DestinationProvider interface {
	CreateDestination(ctx context.Context, input model.Destination) (model.Destination, error)
}

type Response struct {
	Destination *model.Destination  `json:"destination,omitempty"`
	Errors      map[string][]string `json:"errors,omitempty"`
}

type Service struct {
	Provider DestinationProvider
	Logger   *slog.Logger
}

func New(logger *slog.Logger, provider DestinationProvider) Service {
	return Service{Logger: logger, Provider: provider}
}

func (s *Service) CreateDestination(ctx context.Context, input model.DestinationInput) (int, Response) {
	requestId := ctx.Value(utils.RequestIdName)
	response := Response{
		Errors: make(map[string][]string),
	}

	cleanInput, validationErrors := s.validateDestinationInput(input)
	if len(validationErrors) != 0 {
		s.Logger.Debug("Validation errors", "requestId", requestId)
		response.Errors = validationErrors
		return http.StatusBadRequest, response
	}

	destinationTocreate := model.Destination{
		UserId:     ctx.Value(utils.UserId).(int),
		Type:       model.DestinationType(cleanInput.Type),
		Identifier: cleanInput.Identifier,
	}

	newDestination, err := s.Provider.CreateDestination(ctx, destinationTocreate)
	if err != nil {
		s.Logger.Error("Could not create destination from provider", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	response.Destination = &newDestination

	return http.StatusCreated, response
}

func (s *Service) validateDestinationInput(input model.DestinationInput) (model.DestinationInput, map[string][]string) {
	cleanInput, validationErrors := input.Validate()

	return cleanInput, validationErrors
}
