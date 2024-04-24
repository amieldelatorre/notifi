package user

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MockUserProvider struct {
	Users []model.User
}

type TestGetUserByIdTestCase struct {
	ExpectedStatusCode int
	Response           UserResponse
	UserId             int
}

func GetTestUsers() []model.User {
	user1 := model.User{
		Id:              1,
		Email:           "isaac.newton@example.invalid",
		FirstName:       "Isaac",
		LastName:        "Newton",
		Password:        "$argon2id$v=19$m=65536,t=1,p=11$V3dLd/hNnClN0U9mSu3IbQ$mcK9nxoJpUkWsWNbhX14tEC4pXp0oihqJQVePj7FFIc",
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	}

	user2 := model.User{
		Id:              2,
		Email:           "alberteinstein@example.invalid",
		FirstName:       "Albert",
		LastName:        "Einstein",
		Password:        "$argon2id$v=19$m=65536,t=1,p=11$PQT2VdnSXGtAjLKmLHk7jA$hrclADmr/RTFGZgX0J2ujMmZg0adxhOOJczzp1YFMBk",
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	}
	users := []model.User{user1, user2}
	return users
}

type TestDbProviderInstance struct {
	DbPool    *pgxpool.Pool
	Container postgres.PostgresContainer
	Context   context.Context
	Provider  UserProvider
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

	for _, tc := range GetTestUsers() {
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

	provider := repository.NewUserPostgresProvider(dbPool)

	return TestDbProviderInstance{DbPool: dbPool, Container: *postgresContainer, Context: ctx, Provider: provider}
}

func (db *TestDbProviderInstance) CleanUp() {
	// Clean up the container
	if err := db.Container.Terminate(db.Context); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func GetValidTestGetUserByIdTestCases() []TestGetUserByIdTestCase {
	testCases := []TestGetUserByIdTestCase{}

	for _, u := range GetTestUsers() {
		tc := TestGetUserByIdTestCase{
			ExpectedStatusCode: http.StatusOK,
			Response:           UserResponse{User: &u, Errors: map[string][]string{}},
			UserId:             u.Id,
		}
		testCases = append(testCases, tc)
	}
	return testCases
}

func GetInvalidTestGetUserByIdTestCase() TestGetUserByIdTestCase {
	return TestGetUserByIdTestCase{
		ExpectedStatusCode: http.StatusNotFound,
		Response:           UserResponse{Errors: map[string][]string{"user": {"User not found"}}},
		UserId:             100,
	}
}
