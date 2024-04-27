package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/amieldelatorre/notifi/backend/common"
	"github.com/amieldelatorre/notifi/backend/logger"
	"github.com/amieldelatorre/notifi/backend/repository"
	messageDeliveryService "github.com/amieldelatorre/notifi/backend/service/delivery"
	"github.com/amieldelatorre/notifi/backend/utils"
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

	discordProvider := repository.NewDiscordProvider(logger)
	msgProvider := repository.NewMessagePostgresProvider(dbPool)
	destProvider := repository.NewDestinationPostgresProvider(dbPool)
	queueProvider, err := repository.NewSQSMessageQueueProvider(logger, requiredEnvVars.SqsQueueUrl, requiredEnvVars.SqsQueueRegion, requiredEnvVars.SqsQueueName)
	if err != nil {
		logger.Error("Startup failed. Could not connect to the queue", "error", err)
	}

	msgDeliveryService := messageDeliveryService.New(logger, msgProvider, destProvider, &queueProvider, &discordProvider)

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
			time.Sleep(5 * time.Second)
		}
	}
}
