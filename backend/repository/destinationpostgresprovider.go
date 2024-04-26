package repository

import (
	"context"

	"github.com/amieldelatorre/notifi/backend/model"
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

func (p *DestinationPostgresProvider) GetDestinations(ctx context.Context, userId int) ([]model.Destination, error) {
	var destinations []model.Destination

	rows, err := p.DbPool.Query(ctx, `SELECT * FROM Destinations WHERE userId = $1`, userId)
	if err != nil {
		return destinations, err
	}

	for rows.Next() {
		var dest model.Destination

		err := rows.Scan(&dest.Id, &dest.UserId, &dest.Type, &dest.Identifier, &dest.DatetimeCreated, &dest.DatetimeUpdated)
		if err != nil {
			return destinations, err
		}

		destinations = append(destinations, dest)
	}

	if err := rows.Err(); err != nil {
		return destinations, err
	}

	return destinations, nil
}

func (p *DestinationPostgresProvider) GetDestinationById(ctx context.Context, destinationId int, userId int) (model.Destination, error) {
	var destination model.Destination

	// Destination Id's should be unique when querying the destination table
	err := p.DbPool.QueryRow(ctx, "SELECT * FROM Destinations WHERE id = $1 and userId = $2", destinationId, userId).Scan(
		&destination.Id, &destination.UserId, &destination.Type, &destination.Identifier, &destination.DatetimeCreated, &destination.DatetimeUpdated)
	return destination, err
}
