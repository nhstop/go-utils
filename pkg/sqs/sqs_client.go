package sqs

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// NewSQSClient creates and returns an AWS SQS client
func NewSQSClient(ctx context.Context, region string) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	fmt.Println("✅ SQS client initialized")
	return client
}
