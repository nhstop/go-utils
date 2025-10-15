package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// ReceiveConfig holds configurable options for receiving messages
type ReceiveConfig struct {
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	PollInterval        time.Duration // interval to wait on errors
}

// DefaultReceiveConfig provides sensible defaults
func DefaultReceiveConfig() *ReceiveConfig {
	return &ReceiveConfig{
		MaxNumberOfMessages: 5,
		WaitTimeSeconds:     10,
		VisibilityTimeout:   30,
		PollInterval:        5 * time.Second,
	}
}

// MessageHandler is a function type for processing a single message
type MessageHandler func(ctx context.Context, msg types.Message) error

// ReceiveMessages continuously polls SQS and calls handler for each message
func ReceiveMessages(ctx context.Context, client *sqs.Client, queueURL string, cfg *ReceiveConfig, handler MessageHandler) {
	if cfg == nil {
		cfg = DefaultReceiveConfig()
	}

	fmt.Println("👂 Listening for messages on queue:", queueURL)

	for {
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: cfg.MaxNumberOfMessages,
			WaitTimeSeconds:     cfg.WaitTimeSeconds,
			VisibilityTimeout:   cfg.VisibilityTimeout,
		})
		if err != nil {
			fmt.Println("❌ Error receiving messages:", err)
			time.Sleep(cfg.PollInterval)
			continue
		}

		if len(output.Messages) == 0 {
			// No messages, continue polling
			continue
		}

		var wg sync.WaitGroup
		for _, msg := range output.Messages {
			wg.Add(1)
			go func(m types.Message) {
				defer wg.Done()
				if err := handler(ctx, m); err != nil {
					fmt.Println("❌ Error processing message:", err)
				} else {
					// Delete message after successful processing
					_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(queueURL),
						ReceiptHandle: m.ReceiptHandle,
					})
					if err != nil {
						fmt.Println("❌ Failed to delete message:", err)
					} else {
						fmt.Println("✅ Message deleted successfully:", *m.MessageId)
					}
				}
			}(msg)
		}
		wg.Wait()
	}
}
