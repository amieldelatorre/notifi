package model

type QueueMessageBody struct {
	NotifiMessageId int `json:"notifiMessageId"`
}

type QueueMessage struct {
	NotifiMessageId int
	QueueMessageId  string
}
