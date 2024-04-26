package repository

import (
	"encoding/json"
	"fmt"

	"github.com/amieldelatorre/notifi/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSMessageQueueProvider struct {
	Client   sqs.SQS
	QueueUrl string
}

func NewSQSMessageQueueProvider(endpoint string, region string, queueName string) (SQSMessageQueueProvider, error) {
	creds := credentials.AnonymousCredentials
	sessionConfig := &aws.Config{
		Credentials: creds,
	}
	awsSession, err := session.NewSession(sessionConfig)
	if err != nil {
		return SQSMessageQueueProvider{}, err
	}

	sqsConfig := &aws.Config{Endpoint: aws.String(endpoint), Region: aws.String(region)}

	sqsService := sqs.New(awsSession, sqsConfig)
	queueUrl := fmt.Sprintf("%s/queue/%s", endpoint, queueName)

	return SQSMessageQueueProvider{Client: *sqsService, QueueUrl: queueUrl}, nil
}

func (p *SQSMessageQueueProvider) CreateMessage(queueMessageBody model.QueueMessageBody) error {
	messageBody, err := json.Marshal(queueMessageBody)
	if err != nil {
		return err
	}

	messageBodyString := string(messageBody)

	createMessageParams := &sqs.SendMessageInput{
		QueueUrl:    &p.QueueUrl,
		MessageBody: &messageBodyString,
	}

	_, err = p.Client.SendMessage(createMessageParams)
	if err != nil {
		return err
	}

	return nil
}

func (p *SQSMessageQueueProvider) GetMessagesFromQueue(waitTimeSeconds int) ([]model.QueueMessage, error) {
	queueMessages := []model.QueueMessage{}
	receiveMessageParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(p.QueueUrl),
		MaxNumberOfMessages: aws.Int64(3),
		WaitTimeSeconds:     aws.Int64(int64(waitTimeSeconds)),
	}

	response, err := p.Client.ReceiveMessage(receiveMessageParams)

	if err != nil {
		return queueMessages, err
	}

	for _, message := range response.Messages {
		var queueMessageBody model.QueueMessageBody

		err := json.Unmarshal([]byte(*message.Body), &queueMessageBody)
		if err != nil {
			return queueMessages, err
		}

		queueMessage := model.QueueMessage{NotifiMessageId: queueMessageBody.NotifiMessageId, QueueMessageId: *message.ReceiptHandle}
		queueMessages = append(queueMessages, queueMessage)
	}

	return queueMessages, nil
}

func (p *SQSMessageQueueProvider) DeleteMessageFromQueue(id string) error {
	deleteMessageParams := sqs.DeleteMessageInput{
		QueueUrl:      &p.QueueUrl,
		ReceiptHandle: &id,
	}
	_, err := p.Client.DeleteMessage(&deleteMessageParams)
	if err != nil {
		return err
	}

	return nil
}
