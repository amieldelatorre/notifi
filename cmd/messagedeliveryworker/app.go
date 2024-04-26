package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/amieldelatorre/notifi/common"
	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/repository"
	messageDeliveryService "github.com/amieldelatorre/notifi/service/delivery"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	DbPool                 *pgxpool.Pool
	MessageDeliveryService messageDeliveryService.Service
	Logger                 *slog.Logger
}

func NewApp() Application {
	logger := logger.New(os.Stdout, slog.LevelDebug)
	ut := utils.Util{Logger: logger}

	logger.Info("Gathering requirements for application")
	requiredEnvVars, err := ut.GetRequiredEnvVariables()
	if err != nil {
		ut.ExitWithError(1, err)
	}

	st := common.Startup{Logger: logger}
	dbPool := st.InitDb(&requiredEnvVars)

	msgProvider := repository.NewMessagePostgresProvider(dbPool)
	destProvider := repository.NewDestinationPostgresProvider(dbPool)
	queueProvider, err := repository.NewSQSMessageQueueProvider("http://localhost:9324", "ap-southeast2", "notifi")
	if err != nil {
		logger.Error("Startup failed. Could not connect to the queue", "error", err)
	}

	msgDeliveryService := messageDeliveryService.New(logger, msgProvider, destProvider, &queueProvider)

	app := Application{
		DbPool:                 dbPool,
		MessageDeliveryService: msgDeliveryService,
		Logger:                 ut.Logger,
	}
	return app
}

func (app *Application) Run() {
	for {
		if err := app.MessageDeliveryService.ProcessMessages(20); err != nil {
			time.Sleep(5)
		}
	}
}
