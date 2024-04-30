package health

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/backend/logger"
	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/amieldelatorre/notifi/backend/repository"
	"github.com/amieldelatorre/notifi/backend/service/security"
	userService "github.com/amieldelatorre/notifi/backend/service/user"
	"github.com/amieldelatorre/notifi/backend/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDbProviderInstance struct {
	DbPool    *pgxpool.Pool
	Container postgres.PostgresContainer
	Context   context.Context
	Provider  DbProvider
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

		_, err = dbPool.Exec(ctx,
			`INSERT INTO Destinations (userId, type, identifier, datetimeCreated, datetimeUpdated) 
			VALUES ($1, $2, $3, NOW(), NOW())`,
			1, model.DestinationTypeDiscord, "https://one.example.discord.webhook.invalid")
		if err != nil {
			tx.Rollback(ctx)
			panic(err)
		}
	}

	tx.Commit(ctx)

	provider := repository.NewHealthPostgresProvider(dbPool)

	return TestDbProviderInstance{DbPool: dbPool, Container: *postgresContainer, Context: ctx, Provider: &provider}
}

func GetNewMockHealthHandler() (HealthHandler, TestDbProviderInstance, testutils.TestQueueProviderInstance) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := NewTestDbInstance()
	testQueueInstance := testutils.NewTestQueueProviderInstance()

	queueProvider, err := repository.NewSQSMessageQueueProvider(logger, testQueueInstance.Endpoint, "ap-southeast-2", "notifi")
	if err != nil {
		panic(err)
	}
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))

	return HealthHandler{Logger: logger, DbProvider: testDbInstance.Provider, QueueProvider: &queueProvider, JwtService: jwtService}, testDbInstance, testQueueInstance
}

func (db *TestDbProviderInstance) CleanUp() {
	// Clean up the container
	if err := db.Container.Terminate(db.Context); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func TestHealthCheckSuccess(t *testing.T) {
	expectedStatusCode := http.StatusOK
	mockHealthHandler, testDbInstance, testQueueInstance := GetNewMockHealthHandler()
	defer testDbInstance.CleanUp()
	defer testQueueInstance.CleanUp()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()
	mockHealthHandler.healthCheck(response, request)
	result := response.Result()

	if result.StatusCode != expectedStatusCode {
		t.Fatalf("expected status code %d, got %d", expectedStatusCode, result.StatusCode)
	}
}

func TestHealthCheckFail(t *testing.T) {
	expectedStatusCode := http.StatusInternalServerError
	mockHealthHandler, testDbInstance, testQueueInstance := GetNewMockHealthHandler()
	testDbInstance.CleanUp()
	testQueueInstance.CleanUp()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()
	mockHealthHandler.healthCheck(response, request)
	result := response.Result()

	if result.StatusCode != expectedStatusCode {
		t.Fatalf("expected status code %d, got %d", expectedStatusCode, result.StatusCode)
	}
}
