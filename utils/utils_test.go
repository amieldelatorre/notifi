package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGetPostgresConnectionString(t *testing.T) {
	// Expected
	expectedConnectionString := "postgres://root:root@localhost:5432/db"

	// Arrange
	postgresHost := "localhost"
	postgresPort := "5432"
	postgresUsername := "root"
	postgresPassword := "root"
	postgresDabaseName := "db"

	// Act
	actualConnectionString := GetPostgresConnectionString(postgresHost, postgresPort, postgresUsername, postgresPassword, postgresDabaseName)

	// Assert
	if expectedConnectionString != actualConnectionString {
		t.Fatalf("expected: connection string '%s', got '%s'", expectedConnectionString, actualConnectionString)
	}
}

func TestGetRequiredEnvVariableMissing(t *testing.T) {
	fakeEnvVarName := "SHOULD_NOT_EXIST"

	_, err := GetRequiredEnvVariable(fakeEnvVarName)

	if err == nil {
		t.Fatalf("expected error found")
	}
}

func TestGetRequiredEnvVariableFound(t *testing.T) {
	fakeEnvVarName := "SHOULD_EXIST"
	fakeEnvVarValue := "exists"
	os.Setenv(fakeEnvVarName, fakeEnvVarValue)

	value, err := GetRequiredEnvVariable(fakeEnvVarName)

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

	_, err := GetRequiredEnvVariables()

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

	actual, err := GetRequiredEnvVariables()
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

	if os.Getenv("GO_TEST_EXIT_PROGRAM") == "1" {
		ExitWithError(exitStatus, testError)
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
