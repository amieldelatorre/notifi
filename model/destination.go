package model // import "github.com/amieldelatorre/notifi/model"

import (
	"fmt"
	"strings"
	"time"
)

type DestinationType string

const (
	DestinationTypeDiscord = "DISCORD"
)

var DestinationTypes = []DestinationType{DestinationTypeDiscord}

func destinationTypesAsStringSlice() []string {
	strings := []string{}
	for _, dt := range DestinationTypes {
		strings = append(strings, string(dt))
	}
	return strings
}

type Destination struct {
	Id              int             `json:"id"`
	UserId          int             `json:"userId"`
	Type            DestinationType `json:"type"`
	Identifier      string          `json:"identifier"`
	DatetimeCreated time.Time       `json:"datetimeCreated"`
	DatetimeUpdated time.Time       `json:"datetimeUpdated"`
}

type DestinationInput struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func (d *DestinationInput) Validate() (DestinationInput, map[string][]string) {
	validationErrors := make(map[string][]string)

	cleanInput := DestinationInput{
		Type:       strings.ToUpper(strings.TrimSpace(d.Type)),
		Identifier: strings.TrimSpace(d.Identifier),
	}

	userDestinationTypeValid := false
	for _, validDestinationType := range DestinationTypes {
		if cleanInput.Type == string(validDestinationType) {
			userDestinationTypeValid = true
			break
		}
	}
	if !userDestinationTypeValid {
		validationErrors["type"] = append(validationErrors["type"], fmt.Sprintf("Must be one of %s", strings.Join(destinationTypesAsStringSlice(), ", ")))
	}

	if cleanInput.Identifier == "" {
		validationErrors["identifier"] = append(validationErrors["identifier"], "Cannot be empty")
	}

	return cleanInput, validationErrors
}
