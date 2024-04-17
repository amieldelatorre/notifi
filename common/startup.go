package common // import "github.com/amieldelatorre/notifi/common"

import (
	"context"

	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb(requiredEnvVars *utils.RequiredEnvVariables) *pgxpool.Pool {
	postgres_connection_string := utils.GetPostgresConnectionString(
		requiredEnvVars.PortgresHost,
		requiredEnvVars.PortgresPort,
		requiredEnvVars.PortgresUsername,
		requiredEnvVars.PortgresPassword,
		requiredEnvVars.PortgresDabasename,
	)

	dbPool, err := pgxpool.New(context.Background(), postgres_connection_string)
	if err != nil {
		utils.ExitWithError(1, err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		utils.ExitWithError(1, err)
	}

	return dbPool
}
