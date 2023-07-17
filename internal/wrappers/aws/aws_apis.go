package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSAPI represents AWS SDK SQS methods
//
//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_aws/mock_$GOFILE
type SQSAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)

	GetQueueAttributes(ctx context.Context,
		params *sqs.GetQueueAttributesInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// ConfigQueue holds queue params
type ConfigQueue struct {
	QueueName               string `yaml:"name"`
	MaxMessagesPerRetrieval int32  `yaml:"max-messages-per-retrieval"`
	WaitTimeSeconds         int32  `yaml:"wait-time-seconds"`
}
