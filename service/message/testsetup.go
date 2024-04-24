package message // import "github.com/amieldelatorre/notifi/service/message"

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/amieldelatorre/notifi/repository"
	userService "github.com/amieldelatorre/notifi/service/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDbProviderInstance struct {
	DbPool    *pgxpool.Pool
	Container postgres.PostgresContainer
	Context   context.Context
	Provider  MessageProvider
}

func NewTestDbInstance() TestDbProviderInstance {
	ctx := context.Background()

	dbName := "notifi"
	dbUser := "root"
	dbPassword := "root"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres"),
		postgres.WithInitScripts(filepath.Join("../../migrations", "create_tables.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	postgres_connection_string, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	dbPool, err := pgxpool.New(context.Background(), postgres_connection_string)
	if err != nil {
		panic(err)
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		panic(err)
	}

	for _, tc := range userService.GetTestUsers() {
		_, err = dbPool.Exec(ctx,
			`INSERT INTO Users (email, firstName, lastName, password, datetimeCreated, datetimeUpdated) 
			VALUES ($1, $2, $3, $4, NOW(), NOW())`,
			tc.Email, tc.FirstName, tc.LastName, tc.Password)

		if err != nil {
			tx.Rollback(ctx)
			panic(err)
		}
	}
	tx.Commit(ctx)

	provider := repository.NewMessagePostgresProvider(dbPool)

	return TestDbProviderInstance{DbPool: dbPool, Container: *postgresContainer, Context: ctx, Provider: provider}
}

func (db *TestDbProviderInstance) CleanUp() {
	// Clean up the container
	if err := db.Container.Terminate(db.Context); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
