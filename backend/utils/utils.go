package utils // import "github.com/amieldelatorre/notifi/backend/utils"

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const (
	PortgresHostEnvVariableName         = "POSTGRES_HOST"
	PortgresPortEnvVariableName         = "POSTGRES_PORT"
	PortgresUsernameEnvVariableName     = "POSTGRES_USER"
	PortgresPasswordEnvVariableName     = "POSTGRES_PASSWORD"
	PortgresDatabaseNameEnvVariableName = "POSTGRES_DB"
	SqsQueueUrl                         = "SQS_QUEUE_URL"
	SqsQueueRegion                      = "SQS_QUEUE_REGION"
	SqsQueueName                        = "SQS_QUEUE_NAME"
	AwsAccessKeyId                      = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey                  = "AWS_SECRET_ACCESS_KEY"
	AwsSessionToken                     = "AWS_SESSION_TOKEN"
)

type RequiredEnvVariables struct {
	PortgresHost       string
	PortgresPort       string
	PortgresUsername   string
	PortgresPassword   string
	PortgresDabasename string
	SqsQueueUrl        string
	SqsQueueRegion     string
	SqsQueueName       string
}

type OptionalEnvVariables struct {
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	AwsSessionToken    string
}

type Util struct {
	Logger *slog.Logger
}

func (ut *Util) GetRequiredEnvVariables() (RequiredEnvVariables, error) {
	ut.Logger.Debug("Getting required environment variables")
	requiredEnvVariables := RequiredEnvVariables{}
	var retrievalErrors []error

	dbHost, err := ut.GetRequiredEnvVariable(PortgresHostEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbPort, err := ut.GetRequiredEnvVariable(PortgresPortEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbUsername, err := ut.GetRequiredEnvVariable(PortgresUsernameEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbPassword, err := ut.GetRequiredEnvVariable(PortgresPasswordEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbName, err := ut.GetRequiredEnvVariable(PortgresDatabaseNameEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	sqsQueueUrl, err := ut.GetRequiredEnvVariable(SqsQueueUrl)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	sqsQueueRegion, err := ut.GetRequiredEnvVariable(SqsQueueRegion)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	sqsQueueName, err := ut.GetRequiredEnvVariable(SqsQueueName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	if len(retrievalErrors) != 0 {
		ut.Logger.Error("There were missing environment variables")
		return requiredEnvVariables, errors.Join(retrievalErrors...)
	}

	ut.Logger.Debug("All required environment variables found")
	requiredEnvVariables.PortgresHost = dbHost
	requiredEnvVariables.PortgresPort = dbPort
	requiredEnvVariables.PortgresUsername = dbUsername
	requiredEnvVariables.PortgresPassword = dbPassword
	requiredEnvVariables.PortgresDabasename = dbName
	requiredEnvVariables.SqsQueueUrl = sqsQueueUrl
	requiredEnvVariables.SqsQueueRegion = sqsQueueRegion
	requiredEnvVariables.SqsQueueName = sqsQueueName

	return requiredEnvVariables, nil
}

func (ut *Util) GetRequiredEnvVariable(varName string) (string, error) {
	ut.Logger.Debug(fmt.Sprintf("Getting required environment variable '%s'", varName))
	value := strings.TrimSpace(os.Getenv(varName))

	if value == "" {
		ut.Logger.Debug(fmt.Sprintf("Required environment variable '%s'", varName))
		err_msg := fmt.Sprintf("ERROR: environment variable '%s' cannot be blank or empty", varName)
		return "", errors.New(err_msg)
	}

	return value, nil
}

func (ut *Util) ExitWithError(status int, err error) {
	ut.Logger.Error("Exiting with error")
	fmt.Println(err)
	os.Exit(1)
}

func (ut *Util) GetPostgresConnectionString(host string, port string, username string, password string, dbName string) string {
	ut.Logger.Debug("Creating Postgres connection string")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}

func (ut *Util) GetOptionalEnvironmentVariable(varName string) string {
	ut.Logger.Debug(fmt.Sprintf("Getting optional environment variable '%s'", varName))
	value := strings.TrimSpace(os.Getenv(varName))
	return value
}

func (ut *Util) GetOptionalEnvironmentVariables() OptionalEnvVariables {
	optionalEnvVariables := OptionalEnvVariables{
		AwsAccessKeyId:     ut.GetOptionalEnvironmentVariable(AwsAccessKeyId),
		AwsSecretAccessKey: ut.GetOptionalEnvironmentVariable(AwsSecretAccessKey),
		AwsSessionToken:    ut.GetOptionalEnvironmentVariable(AwsSessionToken),
	}

	return optionalEnvVariables
}
