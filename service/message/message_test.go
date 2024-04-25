package message

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"testing"

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/model"
	"github.com/amieldelatorre/notifi/utils"
)

func TestCreateMessage(t *testing.T) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := NewTestDbInstance()
	defer testDbInstance.CleanUp()
	service := New(logger, testDbInstance.Provider, testDbInstance.DestinationProvider)

	tcs := []struct {
		UserId             int
		MessageInput       model.MessageInput
		ExpectedStatusCode int
		ExpectedResponse   Response
	}{
		{
			UserId: 1,
			MessageInput: model.MessageInput{
				Title:         "",
				Body:          "",
				DestinationId: nil,
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: Response{
				Errors: map[string][]string{
					"title":         {"Must have at least one non-whitespace character"},
					"body":          {"Must have at least one non-whitespace character"},
					"destinationId": {"Must be a valid Destination Id"},
				},
			},
		},
		{
			UserId: 1,
			MessageInput: model.MessageInput{
				Title:         "MessageTitle",
				Body:          "MessageBody",
				DestinationId: func(val int) *int { return &val }(1),
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse: Response{
				Errors: map[string][]string{},
				Message: &model.Message{
					Id:     1,
					UserId: 1,
					Title:  "MessageTitle",
					Body:   "MessageBody",
					Status: model.MessageStatusPending,
				},
			},
		},
		{
			UserId: 1,
			MessageInput: model.MessageInput{
				Title:         "MessageTitle",
				Body:          "MessageBody",
				DestinationId: func(val int) *int { return &val }(3),
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: Response{
				Errors: map[string][]string{
					"destinationId": {"Destination Id cannot be found"},
				},
				Message: &model.Message{
					Id:     1,
					UserId: 1,
					Title:  "MessageTitle",
					Body:   "MessageBody",
					Status: model.MessageStatusPending,
				},
			},
		},
		{
			UserId: 2,
			MessageInput: model.MessageInput{
				Title:         "MessageTitle",
				Body:          "MessageBody",
				DestinationId: func(val int) *int { return &val }(1),
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: Response{
				Errors: map[string][]string{
					"destinationId": {"Destination Id cannot be found"},
				},
				Message: &model.Message{
					Id:     1,
					UserId: 2,
					Title:  "MessageTitle",
					Body:   "MessageBody",
					Status: model.MessageStatusPending,
				},
			},
		},
	}

	for tcn, tc := range tcs {
		ctx := context.WithValue(context.Background(), utils.UserId, tc.UserId)
		actualStatusCode, actualResponse := service.CreateMessage(ctx, tc.MessageInput)
		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected response %d, got %d", tcn, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse, actualResponse)
		}

		if actualResponse.Message != nil && tc.ExpectedResponse.Message != nil && (actualResponse.Message.Id != tc.ExpectedResponse.Message.Id ||
			actualResponse.Message.UserId != tc.ExpectedResponse.Message.UserId || actualResponse.Message.Title != tc.ExpectedResponse.Message.Title ||
			actualResponse.Message.Body != tc.ExpectedResponse.Message.Body || actualResponse.Message.Status != tc.ExpectedResponse.Message.Status) {
			t.Fatalf("test case number %d, expected response %+v, got %+v", tcn, tc.ExpectedResponse.Message, actualResponse.Message)
		}
	}
}
