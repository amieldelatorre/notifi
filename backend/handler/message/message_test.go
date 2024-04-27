package message

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

	"github.com/amieldelatorre/notifi/backend/logger"
	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/amieldelatorre/notifi/backend/repository"
	messageService "github.com/amieldelatorre/notifi/backend/service/message"
	"github.com/amieldelatorre/notifi/backend/service/security"
	"github.com/amieldelatorre/notifi/backend/testutils"
	"github.com/amieldelatorre/notifi/backend/utils"
)

func GetNewMockMessageHandler() (MessageHandler, messageService.TestDbProviderInstance, testutils.TestQueueProviderInstance) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := messageService.NewTestDbInstance()

	testQueueInstance := testutils.NewTestQueueProviderInstance()
	queueProvider, err := repository.NewSQSMessageQueueProvider(logger, testQueueInstance.Endpoint, "ap-southeast-2", "notifi")
	if err != nil {
		panic(err)
	}

	msgService := messageService.New(logger, testDbInstance.Provider, testDbInstance.DestinationProvider, &queueProvider)
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))

	mockMessageHandler := New(logger, msgService, jwtService)
	return mockMessageHandler, testDbInstance, testQueueInstance
}

func TestPostMessage(t *testing.T) {
	mockMessageHandler, testDbInstance, testQueueInstance := GetNewMockMessageHandler()
	defer testDbInstance.CleanUp()
	defer testQueueInstance.CleanUp()

	tcs := []struct {
		UserId             int
		MessageInput       model.MessageInput
		ExpectedStatusCode int
		ExpectedResponse   messageService.Response
	}{
		{
			UserId: 1,
			MessageInput: model.MessageInput{
				Title:         "",
				Body:          "",
				DestinationId: nil,
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: messageService.Response{
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
			ExpectedResponse: messageService.Response{
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
			ExpectedResponse: messageService.Response{
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
			ExpectedResponse: messageService.Response{
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
		body, err := json.Marshal(tc.MessageInput)
		if err != nil {
			t.Error(err)
		}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/message", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		ctx := request.Context()
		ctx = context.WithValue(ctx, utils.UserId, tc.UserId)
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		mockMessageHandler.postMessage(response, request)

		result := response.Result()

		if result.StatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case number %d, expected status code %d, got %d", tcn, tc.ExpectedStatusCode, result.StatusCode)
		}

		resultBody, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		result.Body.Close()

		var msgResponse messageService.Response
		err = json.Unmarshal(resultBody, &msgResponse)
		if err != nil {
			t.Error(err)
		}

		if len(msgResponse.Errors) != 0 && len(tc.ExpectedResponse.Errors) != 0 && !reflect.DeepEqual(msgResponse.Errors, tc.ExpectedResponse.Errors) {
			t.Fatalf("test case number %d, expected response errors %+v, got %+v", tcn, tc.ExpectedResponse.Errors, msgResponse.Errors)
		}

		if msgResponse.Message != nil && tc.ExpectedResponse.Message != nil && (msgResponse.Message.Id != tc.ExpectedResponse.Message.Id ||
			msgResponse.Message.UserId != tc.ExpectedResponse.Message.UserId || msgResponse.Message.Title != tc.ExpectedResponse.Message.Title ||
			msgResponse.Message.Body != tc.ExpectedResponse.Message.Body || msgResponse.Message.Status != tc.ExpectedResponse.Message.Status) {
			t.Fatalf("test case number %d, expected response user %+v, got %+v", tcn, tc.ExpectedResponse.Message, msgResponse.Message)
		}

		if tc.ExpectedStatusCode == http.StatusCreated {
			messages, err := mockMessageHandler.Service.QueueProvider.GetMessagesFromQueue(1)
			if err != nil {
				t.Error(err)
			}

			if messages[0].NotifiMessageId != tc.ExpectedResponse.Message.Id {
				t.Fatalf("test case number %d, expected notifi message id %d, got %d", tcn, tc.ExpectedResponse.Message.Id, messages[0].NotifiMessageId)
			}

			err = mockMessageHandler.Service.QueueProvider.DeleteMessageFromQueue(messages[0].QueueMessageId)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
