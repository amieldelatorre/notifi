package model // import "github.com/amieldelatorre/notifi/backend/model"

import (
	"strings"
	"time"
)

type MessageStatus string

const (
	MessageStatusSent    MessageStatus = "SENT"
	MessageStatusPending MessageStatus = "PENDING"
	MessageStatusFailed  MessageStatus = "FAILED"
)

type Message struct {
	Id                  int           `json:"id"`
	UserId              int           `json:"userId"`
	DestinationId       int           `json:"destinationId"`
	Title               string        `json:"title"`
	Body                string        `json:"body"`
	Status              MessageStatus `json:"status"`
	DatetimeCreated     time.Time     `json:"datetimeCreated"`
	DatetimeSendAttempt time.Time     `json:"datetimeSendAttempt"`
}

type MessageInput struct {
	Title         string `json:"title"`
	Body          string `json:"body"`
	DestinationId *int   `json:"destinationId"`
}

func (m *MessageInput) Validate() (MessageInput, map[string][]string) {
	validationErrors := make(map[string][]string)

	cleanInput := MessageInput{
		Title:         m.Title,
		Body:          m.Body,
		DestinationId: m.DestinationId,
	}

	if strings.TrimSpace(m.Title) == "" {
		validationErrors["title"] = append(validationErrors["title"], "Must have at least one non-whitespace character")
	}

	if strings.TrimSpace(m.Body) == "" {
		validationErrors["body"] = append(validationErrors["body"], "Must have at least one non-whitespace character")
	}

	if cleanInput.DestinationId == nil {
		validationErrors["destinationId"] = append(validationErrors["destinationId"], "Must be a valid Destination Id")
	}

	return cleanInput, validationErrors
}
