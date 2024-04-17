package main

import (
	"fmt"
	"net/http"

	"github.com/amieldelatorre/notifi/common"
	"github.com/amieldelatorre/notifi/utils"
	"github.com/jackc/pgx/v5/pgxpool"

	userHandler "github.com/amieldelatorre/notifi/handler/user"
	userProvider "github.com/amieldelatorre/notifi/repository/user"
	userService "github.com/amieldelatorre/notifi/service/user"
)

type Application struct {
	DbPool      *pgxpool.Pool
	UserHandler userHandler.UserHandler
}

func NewApp() Application {
	requiredEnvVars, err := utils.GetRequiredEnvVariables()
	if err != nil {
		utils.ExitWithError(1, err)
	}

	dbPool := common.InitDb(&requiredEnvVars)
	usrHandler := userHandler.New(userService.New(userProvider.NewUserPostgresProvider(dbPool)))

	app := Application{
		DbPool:      dbPool,
		UserHandler: usrHandler,
	}

	return app
}

func (app *Application) Exit() {
	app.DbPool.Close()
}

func (app *Application) Run() {
	mux := http.NewServeMux()
	app.UserHandler.RegisterRoutes(mux)

	fmt.Println("Starting application on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
