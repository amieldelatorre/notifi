package main

import (
	"context"
	"fmt"
	"os"

	UserDb "github.com/amieldelatorre/notifi/internal/db/user"
	UserService "github.com/amieldelatorre/notifi/internal/service/user"
	"github.com/amieldelatorre/notifi/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	requiredEnvVars, err := utils.GetRequiredEnvVariables()
	if err != nil {
		utils.ExitWithError(1, err)
	}

	postgres_connection_string := utils.GetPostgresConnectionString(
		requiredEnvVars.PortgresHost,
		requiredEnvVars.PortgresPort,
		requiredEnvVars.PortgresUsername,
		requiredEnvVars.PortgresPassword,
		requiredEnvVars.PortgresDabasename,
	)

	dbPool, err := pgxpool.New(context.Background(), postgres_connection_string)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	userDb := UserDb.New(dbPool)
	userService := UserService.New(&userDb)

	user, _ := userService.GetUserById(1)

	fmt.Println(user)

	// var user user.User
	// // ctx := context.Background()
	// // transaction, err := dbPool.Begin(ctx)
	// if err != nil {

	// }

	// err = dbPool.QueryRow(context.Background(), "SELECT * FROM users").Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.DatetimeCreated, &user.DatetimeUpdated)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(user)

	// for rows.Next() {
	// 	values, err := rows.Values()
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error iterating through results: %v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	tm := values[5]

	// 	fmt.Println(values[0], values[1], values[2], values[3], values[4], values[5], values[6])
	// }

	// mux := http.NewServeMux()
	// userHandler.RegisterRoutes(mux)

	// fmt.Println("Starting application on port 8080")
	// err := http.ListenAndServe(":8080", mux)
	// if err != nil {
	// 	panic(err)
	// }
}
