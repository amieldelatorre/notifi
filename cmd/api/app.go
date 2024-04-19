package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/amieldelatorre/notifi/common"
	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5/pgxpool"

	AuthHandler "github.com/amieldelatorre/notifi/handler/auth"
	userHandler "github.com/amieldelatorre/notifi/handler/user"
	userProvider "github.com/amieldelatorre/notifi/repository/user"
	AuthService "github.com/amieldelatorre/notifi/service/auth"
	userService "github.com/amieldelatorre/notifi/service/user"
)

type Application struct {
	DbPool      *pgxpool.Pool
	UserHandler userHandler.UserHandler
	AuthHandler AuthHandler.AuthHandler
	Logger      *slog.Logger
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

	usrProvider := userProvider.NewUserPostgresProvider(dbPool)
	usrHandler := userHandler.New(logger, userService.New(logger, usrProvider))
	authHandler := AuthHandler.New(logger, AuthService.New(logger, usrProvider))

	app := Application{
		DbPool:      dbPool,
		UserHandler: usrHandler,
		AuthHandler: authHandler,
		Logger:      logger,
	}

	return app
}

func (app *Application) Exit() {
	app.Logger.Info("Exiting application...")
	app.DbPool.Close()
}

func (app *Application) Run() {
	app.Logger.Info("Attempting to start application...")
	mux := http.NewServeMux()

	app.UserHandler.RegisterRoutes(mux)
	app.AuthHandler.RegisterRoutes(mux)

	app.Logger.Info("Starting application on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
