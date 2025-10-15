package queues

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nhstop/go-utils/pkg/logger"
)

// NewSQSClient creates and returns an AWS SQS client
func NewSQSClient(ctx context.Context, region string) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.Error("failed to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	logger.Info("✅ SQS client initialized")
	return client
}

func SendMessage(ctx context.Context, client *sqs.Client, queueURL, messageBody string) {
	_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		logger.Error("failed to send message: %v", err)
	}

	logger.Info("✅ Message sent successfully")
}
