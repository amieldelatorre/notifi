package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amieldelatorre/notifi/common"
	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5/pgxpool"

	AuthHandler "github.com/amieldelatorre/notifi/handler/auth"
	userHandler "github.com/amieldelatorre/notifi/handler/user"
	userProvider "github.com/amieldelatorre/notifi/repository/user"
	AuthService "github.com/amieldelatorre/notifi/service/auth"
	"github.com/amieldelatorre/notifi/service/security"
	userService "github.com/amieldelatorre/notifi/service/user"
)

type Application struct {
	DbPool      *pgxpool.Pool
	UserHandler userHandler.UserHandler
	AuthHandler AuthHandler.AuthHandler
	Logger      *slog.Logger
	Server      *http.Server
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

	// TODO: Get signing key from environment variable
	signingKey := []byte("super_secret_signing_key")
	jwtService := security.NewJwtService(signingKey)

	usrProvider := userProvider.NewUserPostgresProvider(dbPool)
	usrHandler := userHandler.New(logger, userService.New(logger, usrProvider), jwtService)
	authHandler := AuthHandler.New(logger, AuthService.New(logger, usrProvider, jwtService), jwtService)

	mux := http.NewServeMux()
	usrHandler.RegisterRoutes(mux)
	authHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	app := Application{
		DbPool:      dbPool,
		UserHandler: usrHandler,
		AuthHandler: authHandler,
		Logger:      logger,
		Server:      server,
	}

	return app
}

func (app *Application) Exit() {
	app.Logger.Info("Exiting application...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	app.Logger.Info("Shutting down server")
	err := app.Server.Shutdown(ctx)
	if err != nil {
		app.Logger.Error("Error shutting down server", "error", err)
	}

	app.Logger.Info("Closing database connection")
	app.DbPool.Close()

	app.Logger.Info("Application has been shutdown, bye bye !")
}

func (app *Application) Run() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		app.Logger.Info("Attempting to start application...")
		app.Logger.Info(fmt.Sprintf("Starting application on port %s", app.Server.Addr))
		err := app.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("Something went wrong with the server", "error", err)
		}
	}()

	sig := <-stopChan

	app.Logger.Info(fmt.Sprintf("Received signal '%+v', attempting to shutdown", sig))
	app.Exit()
}
