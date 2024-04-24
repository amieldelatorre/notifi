package message // import "github.com/amieldelatorre/notifi/service/message"

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/utils"
)

type MessageProvider interface {
	CreateMessage(ctx context.Context, input model.Message) (model.Message, error)
}

type Response struct {
	Message *model.Message      `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type Service struct {
	Provider MessageProvider
	Logger   *slog.Logger
}

func New(logger *slog.Logger, provider MessageProvider) Service {
	return Service{Logger: logger, Provider: provider}
}

func (s *Service) CreateMessage(ctx context.Context, input model.MessageInput) (int, Response) {
	requestId := ctx.Value(utils.RequestIdName)
	response := Response{
		Errors: make(map[string][]string),
	}

	cleanInput, validationErrors := s.validateMessageInput(input)
	if len(validationErrors) != 0 {
		s.Logger.Debug("Validation errors", "requestId", requestId)
		response.Errors = validationErrors
		return http.StatusBadRequest, response
	}

	messageToCreate := model.Message{
		UserId: ctx.Value(utils.UserId).(int),
		Title:  cleanInput.Title,
		Body:   cleanInput.Body,
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

func (s *Service) validateMessageInput(input model.MessageInput) (model.MessageInput, map[string][]string) {
	cleanInput, validationErrors := input.Validate()

	return cleanInput, validationErrors
}
