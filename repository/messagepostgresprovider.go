package repository

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagePostgresProvider struct {
	DbPool *pgxpool.Pool
}

func NewMessagePostgresProvider(dbPool *pgxpool.Pool) *MessagePostgresProvider {
	return &MessagePostgresProvider{DbPool: dbPool}
}

func (p *MessagePostgresProvider) CreateMessage(ctx context.Context, input model.Message) (model.Message, error) {
	var newMessage model.Message
	tx, err := p.DbPool.Begin(ctx)
	if err != nil {
		return newMessage, err
	}
	defer tx.Rollback(ctx)

	err = p.DbPool.QueryRow(ctx,
		`INSERT INTO Messages (userId, destinationId, title, body, status, datetimeCreated, datetimeSendAttempt) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
		RETURNING id, destinationId, userId, title, body, status, datetimeCreated, datetimeSendAttempt`,
		input.UserId, input.DestinationId, input.Title, input.Body, model.MessageStatusPending).Scan(
		&newMessage.Id, &newMessage.DestinationId, &newMessage.UserId, &newMessage.Title, &newMessage.Body, &newMessage.Status,
		&newMessage.DatetimeCreated, &newMessage.DatetimeSendAttempt)
	if err != nil {
		return newMessage, err
	}

	err = tx.Commit(ctx)
	return newMessage, err
}

func (p *MessagePostgresProvider) GetMessageById(ctx context.Context, messageId int) (model.Message, error) {
	var message model.Message

	// Destination Id's should be unique when querying the destination table
	err := p.DbPool.QueryRow(ctx, "SELECT * FROM Messages WHERE id = $1", messageId).Scan(
		&message.Id, &message.UserId, &message.DestinationId, &message.Title, &message.Body, &message.Status, &message.DatetimeCreated, &message.DatetimeSendAttempt)
	return message, err
}

func (p *MessagePostgresProvider) UpdateMessage(ctx context.Context, input model.Message) (model.Message, error) {
	var updatedMessage model.Message
	tx, err := p.DbPool.Begin(ctx)
	if err != nil {
		return updatedMessage, err
	}
	defer tx.Rollback(ctx)

	err = p.DbPool.QueryRow(ctx,
		`UPDATE Messages SET
		userId = $1, destinationId = $2, title = $3, body = $4, status = $5, datetimeCreated = $6, datetimeSendAttempt = $7
		WHERE id = $8
		RETURNING id, destinationId, userId, title, body, status, datetimeCreated, datetimeSendAttempt`,
		input.UserId, input.DestinationId, input.Title, input.Body, input.Status, input.DatetimeCreated, input.DatetimeSendAttempt, input.Id).Scan(
		&updatedMessage.Id, &updatedMessage.DestinationId, &updatedMessage.UserId, &updatedMessage.Title, &updatedMessage.Body, &updatedMessage.Status,
		&updatedMessage.DatetimeCreated, &updatedMessage.DatetimeSendAttempt)
	if err != nil {
		return updatedMessage, err
	}

	err = tx.Commit(ctx)
	return updatedMessage, err
}
