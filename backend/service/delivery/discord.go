package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type DestinationDeliveryResponse struct {
	IsApplicationError           bool
	IsDestinationIdentifierError bool
	IsMessageDataError           bool
	Errors                       map[string][]string
}

func (s *Service) DeliverToDiscordWebhook(discordWebhookUrl string, title string, body string) error {

	data := struct {
		Content *string `json:"content"`
		Embeds  []struct {
			Title       string  `json:"title"`
			Description string  `json:"description"`
			Color       *string `json:"color"`
		} `json:"embeds"`
	}{
		Content: nil,
		Embeds: []struct {
			Title       string  `json:"title"`
			Description string  `json:"description"`
			Color       *string `json:"color"`
		}{
			{
				Title:       title,
				Description: body,
				Color:       nil,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		s.Logger.Error("Could not marshal discord data to json", "error", err)
		return err
	}

	response, err := http.Post(discordWebhookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.Logger.Error("Error sending request to discord", "error", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		if response.StatusCode == http.StatusBadRequest {
			// TODO: Add an error column and show user column
			return errors.New("format of message was not valid")
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			s.Logger.Error("Response body for error could not be read from response", "error", err)
			return err
		}

		s.Logger.Error("Message could not be delivered", "statusCode", response.StatusCode, "response", string(body))
		return errors.New("message could not be delivered")
	}

	return nil
}
