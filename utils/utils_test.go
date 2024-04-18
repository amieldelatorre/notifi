package utils

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func GetUtil() Util {
	logger := GetLogger(io.Discard, slog.LevelWarn)
	return Util{Logger: logger}
}

func TestGetLogger(t *testing.T) {
	expectedLoggerType := reflect.TypeOf(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	logger := GetLogger(os.Stdout, slog.LevelWarn)
	loggerType := reflect.TypeOf(logger)

	if loggerType != expectedLoggerType {
		t.Fatalf("expected logger type '%s', got '%s'", expectedLoggerType, loggerType)
	}
}

func TestGetPostgresConnectionString(t *testing.T) {
	// Expected
	expectedConnectionString := "postgres://root:root@localhost:5432/db"

	// Arrange
	postgresHost := "localhost"
	postgresPort := "5432"
	postgresUsername := "root"
	postgresPassword := "root"
	postgresDabaseName := "db"
	ut := GetUtil()

	// Act
	actualConnectionString := ut.GetPostgresConnectionString(postgresHost, postgresPort, postgresUsername, postgresPassword, postgresDabaseName)

	// Assert
	if expectedConnectionString != actualConnectionString {
		t.Fatalf("expected: connection string '%s', got '%s'", expectedConnectionString, actualConnectionString)
	}
}

func TestGetRequiredEnvVariableMissing(t *testing.T) {
	fakeEnvVarName := "SHOULD_NOT_EXIST"
	ut := GetUtil()

	_, err := ut.GetRequiredEnvVariable(fakeEnvVarName)

	if err == nil {
		t.Fatalf("expected error found")
	}
}

func TestGetRequiredEnvVariableFound(t *testing.T) {
	fakeEnvVarName := "SHOULD_EXIST"
	fakeEnvVarValue := "exists"
	os.Setenv(fakeEnvVarName, fakeEnvVarValue)
	ut := GetUtil()

	value, err := ut.GetRequiredEnvVariable(fakeEnvVarName)

	if value != fakeEnvVarValue {
		t.Fatalf("expected environment variable value '%s', got '%s'", fakeEnvVarValue, value)
	}

	if err != nil {
		t.Fatalf("expected no error")
	}
}

func TestGetRequiredEnvVariablesFail(t *testing.T) {
	envVars := []string{
		PortgresHostEnvVariableName,
		PortgresPortEnvVariableName,
		PortgresUsernameEnvVariableName,
		PortgresPasswordEnvVariableName,
		PortgresDatabaseNameEnvVariableName,
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	ut := GetUtil()

	_, err := ut.GetRequiredEnvVariables()

	if err == nil {
		t.Fatalf("expected errors")
	}

	for _, envVar := range envVars {
		if !strings.Contains(err.Error(), fmt.Sprintf("ERROR: environment variable '%s' cannot be blank or empty", envVar)) {
			t.Fatalf("expected '%s' not in error string", envVar)
		}
	}
}

func TestGetRequiredEnvVariablesSuccess(t *testing.T) {
	varValue := "test_value"
	envVars := []string{
		PortgresHostEnvVariableName,
		PortgresPortEnvVariableName,
		PortgresUsernameEnvVariableName,
		PortgresPasswordEnvVariableName,
		PortgresDatabaseNameEnvVariableName,
	}

	for _, envVar := range envVars {
		os.Setenv(envVar, varValue)
	}

	expected := RequiredEnvVariables{
		PortgresHost:       varValue,
		PortgresPort:       varValue,
		PortgresUsername:   varValue,
		PortgresPassword:   varValue,
		PortgresDabasename: varValue,
	}
	ut := GetUtil()

	actual, err := ut.GetRequiredEnvVariables()
	if err != nil {
		t.Fatalf("expected no errors")
	}

	if actual != expected {
		t.Fatalf("expected environment variables '%s', got '%s'", expected, actual)
	}
}

func TestExitWithError(t *testing.T) {
	exitStatus := 1
	testError := errors.New("Test Error")
	ut := GetUtil()

	if os.Getenv("GO_TEST_EXIT_PROGRAM") == "1" {
		ut.ExitWithError(exitStatus, testError)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExitWithError")
	cmd.Env = append(os.Environ(), "GO_TEST_EXIT_PROGRAM=1")
	err := cmd.Run()

	resultingError, typeAssertionOk := err.(*exec.ExitError)

	if !typeAssertionOk {
		t.Fatalf("expected exit with type 'exec.ExitError'")
	}

	if resultingError.Success() {
		t.Fatalf("expected exit code %d, got %d", exitStatus, resultingError.ExitCode())
	}
}
