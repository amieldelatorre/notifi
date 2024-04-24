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
		`INSERT INTO Messages (userId, title, body, status, datetimeCreated, datetimeSendAttempt) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, userId, title, body, status, datetimeCreated, datetimeSendAttempt`,
		input.UserId, input.Title, input.Body, model.MessageStatusPending).Scan(
		&newMessage.Id, &newMessage.UserId, &newMessage.Title, &newMessage.Body, &newMessage.Status,
		&newMessage.DatetimeCreated, &newMessage.DatetimeSendAttempt)
	if err != nil {
		return newMessage, err
	}

	err = tx.Commit(ctx)
	return newMessage, err
}
