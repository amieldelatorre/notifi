package common // import "github.com/amieldelatorre/notifi/common"

import (
	"context"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Startup struct {
	Logger *logger.Logger
}

func (st *Startup) InitDb(requiredEnvVars *utils.RequiredEnvVariables) *pgxpool.Pool {
	st.Logger.Info("Initialising database connection")

	ut := utils.Util{Logger: st.Logger}
	postgres_connection_string := ut.GetPostgresConnectionString(
		requiredEnvVars.PortgresHost,
		requiredEnvVars.PortgresPort,
		requiredEnvVars.PortgresUsername,
		requiredEnvVars.PortgresPassword,
		requiredEnvVars.PortgresDabasename,
	)

	ut.Logger.Debug("Creating database pool with connection string")
	dbPool, err := pgxpool.New(context.Background(), postgres_connection_string)
	if err != nil {
		ut.Logger.Debug("Error when creating database pool")
		ut.ExitWithError(1, err)
	}

	ut.Logger.Debug("Checking connectivity with database")
	err = dbPool.Ping(context.Background())
	if err != nil {
		ut.Logger.Debug("Error with database connectivity")
		ut.ExitWithError(1, err)
	}

	return dbPool
}
