package delivery // import "github.com/amieldelatorre/notifi/backend/service/delivery"

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/jackc/pgx/v5"
)

type QueueProvider interface {
	CreateMessage(queueMessageBody model.QueueMessageBody) error
	GetMessagesFromQueue(waitTimeSeconds int) ([]model.QueueMessage, error)
	DeleteMessageFromQueue(id string) error
}

type MessageProvider interface {
	CreateMessage(ctx context.Context, input model.Message) (model.Message, error)
	GetMessageById(ctx context.Context, messageId int) (model.Message, error)
	UpdateMessage(ctx context.Context, input model.Message) (model.Message, error)
}

type DestinationProvider interface {
	GetDestinationById(ctx context.Context, destinationId int, userId int) (model.Destination, error)
}

type Service struct {
	MessageProvider     MessageProvider
	Logger              *slog.Logger
	DestinationProvider DestinationProvider
	QueueProvider       QueueProvider
}

func New(logger *slog.Logger, messageProvider MessageProvider, destinationProvider DestinationProvider, queueProver QueueProvider) Service {
	return Service{Logger: logger, MessageProvider: messageProvider, DestinationProvider: destinationProvider, QueueProvider: queueProver}
}

func (s *Service) ProcessMessages(waitTimeSeconds int) error {
	queueItems, err := s.QueueProvider.GetMessagesFromQueue(waitTimeSeconds)
	if err != nil {
		s.Logger.Error("Error getting messages from queue", "error", err)
		return err
	}

	for _, item := range queueItems {
		ctx := context.Background()
		messageToDeliver, err := s.MessageProvider.GetMessageById(ctx, item.NotifiMessageId)
		if err != nil {
			s.Logger.Error("Error getting message from database", "error", err)
			return err
		}

		if messageToDeliver.Status != model.MessageStatusPending {
			s.Logger.Warn("Message was already processed", "messageId", messageToDeliver.Id)
			continue
		}

		destination, err := s.DestinationProvider.GetDestinationById(ctx, messageToDeliver.DestinationId, messageToDeliver.UserId)
		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			s.Logger.Error("Error destination for message does not exist")

			updatedMessage := messageToDeliver
			updatedMessage.DatetimeSendAttempt = time.Now().UTC()
			updatedMessage.Status = model.MessageStatusFailed

			_, err = s.MessageProvider.UpdateMessage(ctx, updatedMessage)
			if err != nil {
				s.Logger.Error("Error updating message in database for failed delivery attempt")
				return err
			}

			err = s.QueueProvider.DeleteMessageFromQueue(item.QueueMessageId)
			if err != nil {
				s.Logger.Error("Error removing item from queue")
				return err
			}

			continue

		} else if err != nil {
			s.Logger.Error("Could not get destination from database", "error", err)
			return err
		}

		// errors are already handled in the function itself
		_ = s.DeliverMessage(ctx, destination, messageToDeliver, item)

	}

	return nil
}

func (s *Service) DeliverMessage(ctx context.Context, destination model.Destination, messageToDeliver model.Message, queueItem model.QueueMessage) error {
	updatedMessage := messageToDeliver
	switch destination.Type {
	case model.DestinationTypeDiscord:

		err := s.DeliverToDiscordWebhook(destination.Identifier, messageToDeliver.Title, messageToDeliver.Body)
		if err == nil {
			updatedMessage.DatetimeSendAttempt = time.Now().UTC()
			updatedMessage.Status = model.MessageStatusSent
			s.Logger.Info("Message delivered successfully", "messageId", messageToDeliver.Id)
		} else {
			s.Logger.Debug("Error sending message, updating message status in database")

			updatedMessage.DatetimeSendAttempt = time.Now().UTC()
			updatedMessage.Status = model.MessageStatusFailed
		}
	default:
		s.Logger.Error(fmt.Sprintf("Unknown destination type: %s", destination.Type))
		updatedMessage := messageToDeliver
		updatedMessage.DatetimeSendAttempt = time.Now().UTC()
		updatedMessage.Status = model.MessageStatusFailed
	}

	_, err := s.MessageProvider.UpdateMessage(ctx, updatedMessage)
	if err != nil {
		s.Logger.Error("Error updating message in database for failed delivery attempt", "error", err)
		return err
	}

	err = s.QueueProvider.DeleteMessageFromQueue(queueItem.QueueMessageId)
	if err != nil {
		s.Logger.Error("Error removing item from queue", "error", err)
		return err
	}

	return nil
}
