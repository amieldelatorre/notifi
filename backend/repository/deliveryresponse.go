package repository

type DestinationDeliveryResponse struct {
	IsApplicationError           bool
	IsDestinationIdentifierError bool
	IsMessageDataError           bool
	Errors                       map[string][]string
}
