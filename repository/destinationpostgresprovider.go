package repository

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DestinationPostgresProvider struct {
	DbPool *pgxpool.Pool
}

func NewDestinationPostgresProvider(dbPool *pgxpool.Pool) *DestinationPostgresProvider {
	return &DestinationPostgresProvider{DbPool: dbPool}
}

func (p *DestinationPostgresProvider) CreateDestination(ctx context.Context, input model.Destination) (model.Destination, error) {
	var newDestination model.Destination
	tx, err := p.DbPool.Begin(ctx)
	if err != nil {
		return newDestination, err
	}
	defer tx.Rollback(ctx)

	err = p.DbPool.QueryRow(ctx,
		`INSERT INTO Destinations (userId, type, identifier, datetimeCreated, datetimeUpdated) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id, userId, type, identifier, datetimeCreated, datetimeUpdated`,
		input.UserId, input.Type, input.Identifier).Scan(
		&newDestination.Id, &newDestination.UserId, &newDestination.Type, &newDestination.Identifier,
		&newDestination.DatetimeCreated, &newDestination.DatetimeUpdated)
	if err != nil {
		return newDestination, err
	}

	err = tx.Commit(ctx)
	return newDestination, err
}
