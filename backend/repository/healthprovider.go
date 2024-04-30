package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthPostgresProvider struct {
	DbPool *pgxpool.Pool
}

func NewHealthPostgresProvider(dbPool *pgxpool.Pool) HealthPostgresProvider {
	return HealthPostgresProvider{DbPool: dbPool}
}

func (p *HealthPostgresProvider) IsHealthy(ctx context.Context) bool {
	err := p.DbPool.Ping(context.Background())
	return err == nil
}
