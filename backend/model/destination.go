package model // import "github.com/amieldelatorre/notifi/backend/model"

import (
	"fmt"
	"strings"
	"time"
)

type DestinationType string

const (
	DestinationTypeDiscord       = "DISCORD"
	DestinationTypeMobileAndroid = "MOBILE_ANDROID"
	DestinationTypeMobileIos     = "MOBILE_IOS"
)

var DestinationTypes = []DestinationType{DestinationTypeDiscord, DestinationTypeMobileAndroid, DestinationTypeMobileIos}

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

	userDestinationTypeValid := validDestinationType(cleanInput.Type)
	if !userDestinationTypeValid {
		validationErrors["type"] = append(validationErrors["type"], fmt.Sprintf("Must be one of %s", strings.Join(destinationTypesAsStringSlice(), ", ")))
	}

	if cleanInput.Identifier == "" {
		validationErrors["identifier"] = append(validationErrors["identifier"], "Cannot be empty")
	}

	return cleanInput, validationErrors
}

func validDestinationType(destinationType string) bool {
	switch destinationType {
	case DestinationTypeDiscord, DestinationTypeMobileIos, DestinationTypeMobileAndroid:
		return true
	default:
		return false
	}
}
