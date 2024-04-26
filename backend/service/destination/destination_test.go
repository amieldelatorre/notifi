package destination

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/backend/logger"
	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/amieldelatorre/notifi/backend/utils"
)

func TestCreateDestination(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := NewTestDbInstance()
	defer testDbInstance.CleanUp()
	service := New(logger, testDbInstance.Provider)

	tcs := []struct {
		DestinationInput   model.DestinationInput
		ExpectedStatusCode int
		ExpectedResponse   Response
	}{
		{
			DestinationInput: model.DestinationInput{
				Type:       "",
				Identifier: "",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: Response{
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
			ExpectedResponse: Response{
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
			ExpectedResponse: Response{
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
		},
	}

	for tcn, tc := range tcs {
		ctx := context.WithValue(context.Background(), utils.UserId, 1)
		actualStatusCode, actualResponse := service.CreateDestination(ctx, tc.DestinationInput)
		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse, actualResponse)
		}

		if actualResponse.Destination != nil && tc.ExpectedResponse.Destination != nil && (actualResponse.Destination.Id != tc.ExpectedResponse.Destination.Id ||
			actualResponse.Destination.UserId != tc.ExpectedResponse.Destination.UserId || actualResponse.Destination.Type != tc.ExpectedResponse.Destination.Type ||
			actualResponse.Destination.Identifier != tc.ExpectedResponse.Destination.Identifier) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse.Destination, actualResponse.Destination)
		}
	}
}

func TestGetAllDestinations(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := NewTestDbInstance()
	defer testDbInstance.CleanUp()
	service := New(logger, testDbInstance.Provider)

	tcs := []struct {
		UserId             int
		ExpectedStatusCode int
		ExpectedResponse   GetAllResponse
	}{
		{
			UserId:             2,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: GetAllResponse{
				Destinations: []model.Destination{},
				Errors:       map[string][]string{},
			},
		},
		{
			UserId:             1,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: GetAllResponse{
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
				Errors: map[string][]string{},
			},
		},
	}

	for tcn, tc := range tcs {
		ctx := context.WithValue(context.Background(), utils.UserId, 1)
		actualStatusCode, actualResponse := service.GetAllDestinations(ctx, tc.UserId)
		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse, actualResponse)
		}

		if len(tc.ExpectedResponse.Destinations) != len(actualResponse.Destinations) {
			t.Fatalf("test case number %d, expected number of destinations %d, got %d", tcn, len(tc.ExpectedResponse.Destinations), len(actualResponse.Destinations))
		}

		for index, expectedDests := range tc.ExpectedResponse.Destinations {
			responseDests := actualResponse.Destinations[index]

			if responseDests.Id != expectedDests.Id || responseDests.UserId != expectedDests.UserId ||
				responseDests.Type != expectedDests.Type ||
				responseDests.Identifier != expectedDests.Identifier {
				t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse.Destinations[index], responseDests)
			}
		}
	}
}

func TestGetDestinationById(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := NewTestDbInstance()
	defer testDbInstance.CleanUp()
	service := New(logger, testDbInstance.Provider)

	tcs := []struct {
		DestinationId      int
		UserId             int
		ExpectedStatusCode int
		ExpectedResponse   Response
	}{
		{
			DestinationId:      1,
			UserId:             2,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse: Response{
				Destination: nil,
				Errors: map[string][]string{
					"destination": {"Destination not found"},
				},
			},
		},
		{
			DestinationId:      3,
			UserId:             1,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse: Response{
				Destination: nil,
				Errors: map[string][]string{
					"destination": {"Destination not found"},
				},
			},
		},
		{
			DestinationId:      1,
			UserId:             1,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: Response{
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
		ctx := context.Background()
		actualStatusCode, actualResponse := service.GetDestinationById(ctx, tc.DestinationId, tc.UserId)
		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse, actualResponse)
		}

		if actualResponse.Destination != nil && tc.ExpectedResponse.Destination != nil && (actualResponse.Destination.Id != tc.ExpectedResponse.Destination.Id ||
			actualResponse.Destination.UserId != tc.ExpectedResponse.Destination.UserId || actualResponse.Destination.Type != tc.ExpectedResponse.Destination.Type ||
			actualResponse.Destination.Identifier != tc.ExpectedResponse.Destination.Identifier) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse.Destination, actualResponse.Destination)
		}
	}
}
