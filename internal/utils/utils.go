package utils // import "github.com/amieldelatorre/notifi/internal/utils"

import (
	"errors"
	"fmt"
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

func GetRequiredEnvVariables() (RequiredEnvVariables, error) {
	requiredEnvVariables := RequiredEnvVariables{}
	var retrievalErrors []error

	dbHost, err := GetRequiredEnvVariable(PortgresHostEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbPort, err := GetRequiredEnvVariable(PortgresPortEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbUsername, err := GetRequiredEnvVariable(PortgresUsernameEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbPassword, err := GetRequiredEnvVariable(PortgresPasswordEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	dbName, err := GetRequiredEnvVariable(PortgresDatabaseNameEnvVariableName)
	if err != nil {
		retrievalErrors = append(retrievalErrors, err)
	}

	if len(retrievalErrors) != 0 {
		return requiredEnvVariables, errors.Join(retrievalErrors...)
	}

	requiredEnvVariables.PortgresHost = dbHost
	requiredEnvVariables.PortgresPort = dbPort
	requiredEnvVariables.PortgresUsername = dbUsername
	requiredEnvVariables.PortgresPassword = dbPassword
	requiredEnvVariables.PortgresDabasename = dbName

	return requiredEnvVariables, nil
}

func GetRequiredEnvVariable(varName string) (string, error) {
	value := os.Getenv(varName)

	if value == "" {
		err_msg := fmt.Sprintf("ERROR: environment variable '%s' cannot be blank or empty", varName)
		return "", errors.New(err_msg)
	}

	return value, nil
}

func ExitWithError(status int, err error) {
	fmt.Println(err)
	os.Exit(1)
}

func GetPostgresConnectionString(host string, port string, username string, password string, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}
