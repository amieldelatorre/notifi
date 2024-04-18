package utils // import "github.com/amieldelatorre/notifi/utils"

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

const (
	PortgresHostEnvVariableName         = "POSTGRES_HOST"
	PortgresPortEnvVariableName         = "POSTGRES_PORT"
	PortgresUsernameEnvVariableName     = "POSTGRES_USER"
	PortgresPasswordEnvVariableName     = "POSTGRES_PASSWORD"
	PortgresDatabaseNameEnvVariableName = "POSTGRES_DB"
)

type RequiredEnvVariables struct {
	PortgresHost       string
	PortgresPort       string
	PortgresUsername   string
	PortgresPassword   string
	PortgresDabasename string
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

	if len(retrievalErrors) != 0 {
		ut.Logger.Debug("There were missing environment variables")
		return requiredEnvVariables, errors.Join(retrievalErrors...)
	}

	ut.Logger.Debug("All required environment variables found")
	requiredEnvVariables.PortgresHost = dbHost
	requiredEnvVariables.PortgresPort = dbPort
	requiredEnvVariables.PortgresUsername = dbUsername
	requiredEnvVariables.PortgresPassword = dbPassword
	requiredEnvVariables.PortgresDabasename = dbName

	return requiredEnvVariables, nil
}

func (ut *Util) GetRequiredEnvVariable(varName string) (string, error) {
	ut.Logger.Debug("Getting required environment variable '", varName, "'")
	value := os.Getenv(varName)

	if value == "" {
		ut.Logger.Debug("Required environment variable '", varName, "' not found")
		err_msg := fmt.Sprintf("ERROR: environment variable '%s' cannot be blank or empty", varName)
		return "", errors.New(err_msg)
	}

	return value, nil
}

func (ut *Util) ExitWithError(status int, err error) {
	ut.Logger.Debug("Exiting with error")
	fmt.Println(err)
	os.Exit(1)
}

func (ut *Util) GetPostgresConnectionString(host string, port string, username string, password string, dbName string) string {
	ut.Logger.Debug("Creating Postgres connection string")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}
