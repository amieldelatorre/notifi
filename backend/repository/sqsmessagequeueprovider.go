package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/amieldelatorre/notifi/backend/model"
	"github.com/amieldelatorre/notifi/backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSMessageQueueProvider struct {
	Client   sqs.SQS
	QueueUrl string
}

func NewSQSMessageQueueProvider(logger *slog.Logger, endpoint string, region string, queueName string) (SQSMessageQueueProvider, error) {
	logger.Info("Attempting to connect to SQS")
	ut := utils.Util{Logger: logger}
	optionalEnvVariables := ut.GetOptionalEnvironmentVariables()
	var sessionConfig *aws.Config

	if optionalEnvVariables.AwsAccessKeyId != "" || optionalEnvVariables.AwsSecretAccessKey != "" || optionalEnvVariables.AwsSessionToken != "" {
		logger.Info("One of the AWS environment variables is not empty, using given credentials even if the others are empty")
		creds := credentials.NewStaticCredentials(optionalEnvVariables.AwsAccessKeyId, optionalEnvVariables.AwsSecretAccessKey, optionalEnvVariables.AwsSessionToken)

		sessionConfig = &aws.Config{
			Region:      &region,
			Credentials: creds,
		}
	} else {
		logger.Info("No AWS environment variables found, using default provider chain to look for credentials")
		sessionConfig = &aws.Config{
			Region: &region,
		}
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

func (p *SQSMessageQueueProvider) IsHealthy(ctx context.Context) bool {
	input := sqs.GetQueueAttributesInput{
		QueueUrl: &p.QueueUrl,
	}
	_, err := p.Client.GetQueueAttributes(&input)
	return err == nil
}
