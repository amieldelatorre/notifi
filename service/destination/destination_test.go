package destination

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/utils"
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
					Id:              1,
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
