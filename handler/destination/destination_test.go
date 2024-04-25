package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/model"
	destinationService "github.com/amieldelatorre/notifi/service/destination"
	"github.com/amieldelatorre/notifi/service/security"
	"github.com/amieldelatorre/notifi/utils"
)

func GetNewMockDestinationHandler() (DestinationHandler, destinationService.TestDbProviderInstance) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := destinationService.NewTestDbInstance()
	destService := destinationService.New(logger, testDbInstance.Provider)
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))

	mockDestinationHandler := New(logger, destService, jwtService)
	return mockDestinationHandler, testDbInstance
}

func TestPostDestination(t *testing.T) {
	mockDestinationHandler, testDbInstance := GetNewMockDestinationHandler()
	defer testDbInstance.CleanUp()

	tcs := []struct {
		DestinationInput   model.DestinationInput
		ExpectedStatusCode int
		ExpectedResponse   destinationService.Response
	}{{
		DestinationInput: model.DestinationInput{
			Type:       "",
			Identifier: "",
		},
		ExpectedStatusCode: http.StatusBadRequest,
		ExpectedResponse: destinationService.Response{
			Errors: map[string][]string{
				"type":       {"Must be one of DISCORD"},
				"identifier": {"Cannot be empty"},
			},
		},
	},
		{
			DestinationInput: model.DestinationInput{
				Type:       "x",
				Identifier: "anidentifier",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: destinationService.Response{
				Errors: map[string][]string{
					"type": {"Must be one of DISCORD"},
				},
			},
		},
		{
			DestinationInput: model.DestinationInput{
				Type:       "DISCORD ",
				Identifier: "anidentifier",
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse: destinationService.Response{
				Destination: &model.Destination{
					Id:              1,
					UserId:          1,
					Type:            "DISCORD",
					Identifier:      "anidentifier",
					DatetimeCreated: time.Now(),
					DatetimeUpdated: time.Now(),
				},
				Errors: map[string][]string{},
			},
		}}

	for tcn, tc := range tcs {
		body, err := json.Marshal(tc.DestinationInput)
		if err != nil {
			t.Error(err)
		}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/destination", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		ctx := request.Context()
		ctx = context.WithValue(ctx, utils.UserId, 1)
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		mockDestinationHandler.postDestination(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected status code %d, got %d", tcn, tc.ExpectedStatusCode, result.StatusCode)
		}

		resultBody, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		result.Body.Close()

		var destResponse destinationService.Response
		err = json.Unmarshal(resultBody, &destResponse)
		if err != nil {
			t.Error(err)
		}

		if len(destResponse.Errors) != 0 && len(tc.ExpectedResponse.Errors) != 0 && !reflect.DeepEqual(destResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response errors %+v, got %+v", tcn, tc.ExpectedResponse.Errors, destResponse.Errors)
		}

		if destResponse.Destination != nil && tc.ExpectedResponse.Destination != nil && (destResponse.Destination.Id != tc.ExpectedResponse.Destination.Id ||
			destResponse.Destination.UserId != tc.ExpectedResponse.Destination.UserId || destResponse.Destination.Identifier != tc.ExpectedResponse.Destination.Identifier) {
			t.Fatalf("test case number %d, expected response user %+v, got %+v", tcn, tc.ExpectedResponse.Destination, destResponse.Destination)
		}
	}
}
