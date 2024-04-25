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
					Id:              3,
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

func TestGetDestinations(t *testing.T) {
	mockDestinationHandler, testDbInstance := GetNewMockDestinationHandler()
	defer testDbInstance.CleanUp()

	tcs := []struct {
		UserId             int
		ExpectedStatusCode int
		ExpectedResponse   destinationService.GetAllResponse
	}{
		{
			UserId:             2,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: destinationService.GetAllResponse{
				Destinations: []model.Destination{},
			},
		},
		{
			UserId:             1,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: destinationService.GetAllResponse{
				Destinations: []model.Destination{
					{
						Id:              1,
						UserId:          1,
						Type:            model.DestinationTypeDiscord,
						Identifier:      "https://one.example.discord.webhook.invalid",
						DatetimeCreated: time.Now(),
						DatetimeUpdated: time.Now(),
					},
					{
						Id:              2,
						UserId:          1,
						Type:            model.DestinationTypeDiscord,
						Identifier:      "https://two.example.discord.webhook.invalid",
						DatetimeCreated: time.Now(),
						DatetimeUpdated: time.Now(),
					},
				},
			},
		},
	}

	for tcn, tc := range tcs {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/destination", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, utils.UserId, tc.UserId)
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		mockDestinationHandler.getDestinations(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, result.StatusCode)
		}

		resultBody, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		result.Body.Close()

		var destResponse destinationService.GetAllResponse
		err = json.Unmarshal(resultBody, &destResponse)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(destResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse, destResponse)
		}

		if len(tc.ExpectedResponse.Destinations) != len(destResponse.Destinations) {
			t.Fatalf("test case number %d, expected number of destinations %d, got %d", tcn, len(tc.ExpectedResponse.Destinations), len(destResponse.Destinations))
		}

		for index, expectedDests := range tc.ExpectedResponse.Destinations {
			responseDests := destResponse.Destinations[index]

			if responseDests.Id != expectedDests.Id || responseDests.UserId != expectedDests.UserId ||
				responseDests.Type != expectedDests.Type ||
				responseDests.Identifier != expectedDests.Identifier {
				t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse.Destinations[index], responseDests)
			}
		}
	}
}

func TestGetDestinationById(t *testing.T) {
	mockDestinationHandler, testDbInstance := GetNewMockDestinationHandler()
	defer testDbInstance.CleanUp()

	tcs := []struct {
		DestinationId      string
		UserId             int
		ExpectedStatusCode int
		ExpectedResponse   destinationService.Response
	}{
		{
			DestinationId:      "1",
			UserId:             2,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse: destinationService.Response{
				Destination: nil,
				Errors: map[string][]string{
					"destination": {"Destination not found"},
				},
			},
		},
		{
			DestinationId:      "3",
			UserId:             1,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse: destinationService.Response{
				Destination: nil,
				Errors: map[string][]string{
					"destination": {"Destination not found"},
				},
			},
		},
		{
			DestinationId:      "invalid",
			UserId:             1,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: destinationService.Response{
				Destination: nil,
				Errors: map[string][]string{
					"id": {"Invalid destination Id provided. Must be an integer"},
				},
			},
		},
		{
			DestinationId:      "1",
			UserId:             1,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: destinationService.Response{
				Destination: &model.Destination{
					Id:              1,
					UserId:          1,
					Type:            model.DestinationTypeDiscord,
					Identifier:      "https://one.example.discord.webhook.invalid",
					DatetimeCreated: time.Now(),
					DatetimeUpdated: time.Now(),
				},
				Errors: map[string][]string{},
			},
		},
	}

	for tcn, tc := range tcs {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/destination", nil)
		request.SetPathValue("id", tc.DestinationId)
		ctx := request.Context()
		ctx = context.WithValue(ctx, utils.UserId, tc.UserId)
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		mockDestinationHandler.getDestinationById(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, result.StatusCode)
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

		// resultBody, err := io.ReadAll(result.Body)
		// if err != nil {
		// 	t.Error(err)
		// }
		// result.Body.Close()

		// ctx := context.Background()
		// actualStatusCode, actualResponse := service.GetDestinationById(ctx, tc.DestinationId, tc.UserId)
		// if actualStatusCode != tc.ExpectedStatusCode {
		// 	t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, actualStatusCode)
		// }

		// if !reflect.DeepEqual(actualResponse.Errors, tc.ExpectedResponse.Errors) {
		// 	t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse, actualResponse)
		// }

		// if actualResponse.Destination != nil && tc.ExpectedResponse.Destination != nil && (actualResponse.Destination.Id != tc.ExpectedResponse.Destination.Id ||
		// 	actualResponse.Destination.UserId != tc.ExpectedResponse.Destination.UserId || actualResponse.Destination.Type != tc.ExpectedResponse.Destination.Type ||
		// 	actualResponse.Destination.Identifier != tc.ExpectedResponse.Destination.Identifier) {
		// 	t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse.Destination, actualResponse.Destination)
		// }
	}
}
