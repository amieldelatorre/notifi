package message // import "github.com/amieldelatorre/notifi/service/message"

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5"
)

type MessageProvider interface {
	CreateMessage(ctx context.Context, input model.Message) (model.Message, error)
}

type DestinationProvider interface {
	GetDestinationById(ctx context.Context, destinationId int, userId int) (model.Destination, error)
}

type Response struct {
	Message *model.Message      `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type Service struct {
	Provider            MessageProvider
	Logger              *slog.Logger
	DestinationProvider DestinationProvider
}

func New(logger *slog.Logger, provider MessageProvider, destinationProvider DestinationProvider) Service {
	return Service{Logger: logger, Provider: provider, DestinationProvider: destinationProvider}
}

func (s *Service) CreateMessage(ctx context.Context, input model.MessageInput) (int, Response) {
	requestId := ctx.Value(utils.RequestIdName)
	response := Response{
		Errors: make(map[string][]string),
	}

	userId := ctx.Value(utils.UserId).(int)

	cleanInput, validationErrors, err := s.validateMessageInput(ctx, input, userId)
	if err != nil {
		s.Logger.Error("Could not validate destinationInput", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	} else if len(validationErrors) != 0 {
		s.Logger.Debug("Validation errors", "requestId", requestId)
		response.Errors = validationErrors
		return http.StatusBadRequest, response
	}

	messageToCreate := model.Message{
		UserId:        userId,
		DestinationId: *cleanInput.DestinationId,
		Title:         cleanInput.Title,
		Body:          cleanInput.Body,
	}

	newMessage, err := s.Provider.CreateMessage(ctx, messageToCreate)
	if err != nil {
		s.Logger.Error("Could not create message from provider", "requestId", requestId, "error", err)
		response.Errors["server"] = append(response.Errors["server"], "Something went wrong")
		return http.StatusInternalServerError, response
	}

	response.Message = &newMessage
	return http.StatusCreated, response
}

func (s *Service) DestinationIdExists(ctx context.Context, destinationId int, userId int) (bool, error) {
	_, err := s.DestinationProvider.GetDestinationById(ctx, destinationId, userId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) validateMessageInput(ctx context.Context, input model.MessageInput, userId int) (model.MessageInput, map[string][]string, error) {
	cleanInput, validationErrors := input.Validate()

	if len(validationErrors) > 0 {
		return cleanInput, validationErrors, nil
	}

	validDestinationId, err := s.DestinationIdExists(ctx, *cleanInput.DestinationId, userId)
	if err != nil {
		return cleanInput, validationErrors, err
	}
	if !validDestinationId {
		validationErrors["destinationId"] = append(validationErrors["destinationId"], "Destination Id cannot be found")
	}

	return cleanInput, validationErrors, nil
}
