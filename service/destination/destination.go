package destination // import "github.com/amieldelatorre/notifi/service/destination"

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5"
)

type DestinationProvider interface {
	CreateDestination(ctx context.Context, input model.Destination) (model.Destination, error)
	GetDestinations(ctx context.Context, userId int) ([]model.Destination, error)
	GetDestinationById(ctx context.Context, destinationId int, userId int) (model.Destination, error)
}

type Response struct {
	Destination *model.Destination  `json:"destination,omitempty"`
	Errors      map[string][]string `json:"errors,omitempty"`
}

type GetAllResponse struct {
	Destinations []model.Destination `json:"destinations"`
	Errors       map[string][]string `json:"errors,omitempty"`
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

func (s *Service) GetAllDestinations(ctx context.Context, userId int) (int, GetAllResponse) {
	requestId := ctx.Value(utils.RequestIdName)
	response := GetAllResponse{
		Errors: make(map[string][]string),
	}

	destinations, err := s.Provider.GetDestinations(ctx, userId)
	if err != nil {
		s.Logger.Error("Could not get destinations from provider", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	response.Destinations = destinations

	return http.StatusOK, response
}

func (s *Service) GetDestinationById(ctx context.Context, destinationId int, userId int) (int, Response) {
	requestId := ctx.Value(utils.RequestIdName)
	response := Response{
		Errors: make(map[string][]string),
	}

	destination, err := s.Provider.GetDestinationById(ctx, destinationId, userId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		response.Errors["destination"] = append(response.Errors["destination"], "Destination not found")
		return http.StatusNotFound, response
	} else if err != nil {
		s.Logger.Error("Could not get destination from provider", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	response.Destination = &destination
	return http.StatusOK, response
}

func (s *Service) validateDestinationInput(input model.DestinationInput) (model.DestinationInput, map[string][]string) {
	cleanInput, validationErrors := input.Validate()

	return cleanInput, validationErrors
}
