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

	"github.com/amieldelatorre/notifi/logger"
	"github.com/amieldelatorre/notifi/model"
	messageService "github.com/amieldelatorre/notifi/service/message"
	"github.com/amieldelatorre/notifi/service/security"
	"github.com/amieldelatorre/notifi/utils"
)

func GetNewMockMessageHandler() (MessageHandler, messageService.TestDbProviderInstance) {
	logger := logger.New(io.Discard, slog.LevelWarn)
	testDbInstance := messageService.NewTestDbInstance()
	msgService := messageService.New(logger, testDbInstance.Provider)
	jwtService := security.NewJwtService([]byte("super_secret_signing_key"))

	mockMessageHandler := New(logger, msgService, jwtService)
	return mockMessageHandler, testDbInstance
}

func TestPostMessage(t *testing.T) {
	mockMessageHandler, testDbInstance := GetNewMockMessageHandler()
	defer testDbInstance.CleanUp()

	tcs := []struct {
		MessageInput       model.MessageInput
		ExpectedStatusCode int
		ExpectedResponse   messageService.Response
	}{
		{
			MessageInput: model.MessageInput{
				Title: "",
				Body:  "",
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse: messageService.Response{
				Errors: map[string][]string{
					"title": {"Must have at least one non-whitespace character"},
					"body":  {"Must have at least one non-whitespace character"},
				},
			},
		},
		{
			MessageInput: model.MessageInput{
				Title: "MessageTitle",
				Body:  "MessageBody",
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
	}

	for tcn, tc := range tcs {
		body, err := json.Marshal(tc.MessageInput)
		if err != nil {
			t.Error(err)
		}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/message", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		ctx := request.Context()
		ctx = context.WithValue(ctx, utils.UserId, 1)
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
	}
}
