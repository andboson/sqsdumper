package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"
)

// MessageHandler represents a single SQS message handler
type MessageHandler func(poller SQSPoller, msg types.Message) error

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_sqs/mock_$GOFILE

// SQSPoller represents a long-polling Amazon SQS queue
type SQSPoller interface {
	GetQueueURL() *string
	GetTotal() int
	PollMessages(ctx context.Context, messageHandler MessageHandler) error
	DeleteMessage(ctx context.Context, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	GetQueueAttrs(ctx context.Context, input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error)
}

type sqsPoller struct {
	client        SQSAPI
	logger        zerolog.Logger
	cfg           ConfigQueue
	queueURL      *string
	stopOnTotal   bool
	stopAfter     int
	totalMessages int
	checkReceived map[string]struct{}
	counterChan   chan int
	bar           *progressbar.ProgressBar
}

// SQSParam holds SQSPoller params
type SQSParam struct {
	Client      SQSAPI
	Logger      zerolog.Logger
	QueueConfig ConfigQueue
	StopOnTotal bool
	StopAfter   int
	CounterChan chan int
}

// NewSQSPoller returns an instance of SQSPoller
func NewSQSPoller(params SQSParam) (SQSPoller, error) {
	s := &sqsPoller{
		client:        params.Client,
		logger:        params.Logger,
		cfg:           params.QueueConfig,
		stopOnTotal:   params.StopOnTotal,
		counterChan:   params.CounterChan,
		stopAfter:     params.StopAfter,
		checkReceived: map[string]struct{}{},
	}
	queueURL, err := s.fetchQueueURL(context.Background(), s.cfg.QueueName)
	if err != nil {
		return nil, errors.Wrap(err, "error getting AWS SQS queue URL")
	}

	s.queueURL = queueURL.QueueUrl

	queueAttrs, err := s.GetQueueAttrs(context.Background(), &sqs.GetQueueAttributesInput{
		QueueUrl: queueURL.QueueUrl,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
		},
	})

	s.totalMessages, err = strconv.Atoi(queueAttrs.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessages)])
	if err != nil {
		return nil, errors.Wrap(err, "error getting total number of messages from AWS SQS queue URL")
	}

	s.bar = progressbar.Default(int64(s.totalMessages), "Processing..")

	return s, nil
}

func (s *sqsPoller) GetQueueURL() *string {
	return s.queueURL
}

func (s *sqsPoller) GetTotal() int {
	return s.totalMessages
}

func (s *sqsPoller) PollMessages(ctx context.Context, messageHandler MessageHandler) error {
	if messageHandler == nil {
		return errors.New("a message handler is nil, stopped")
	}
	var processed int

	// start polling
	for {
		select {
		case <-ctx.Done():
			s.logger.Log().Msg("got context.Done signal, exiting processing")
			return nil
		default:
			output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            s.queueURL,
				MaxNumberOfMessages: s.cfg.MaxMessagesPerRetrieval,
				WaitTimeSeconds:     s.cfg.WaitTimeSeconds,
			})
			if err != nil {
				s.logger.Err(err).Msg("can't get new messages from SQS")
				continue
			}
			for _, message := range output.Messages {
				if err := messageHandler(s, message); err != nil {
					s.logger.Err(err).Msg("processing error")
				}
				s.bar.Add(1)

				if s.stopAfter != 0 && processed >= s.stopAfter {
					fmt.Printf("\n")
					s.logger.Log().Msgf("stopped after %d messages processed", processed)
					return nil
				}

				if s.bar.IsFinished() && s.stopOnTotal {
					fmt.Printf("\n")
					s.logger.Log().Msg("all messages processed")
					return nil
				}
			}
		}
	}

	return nil
}

func (s *sqsPoller) fetchQueueURL(ctx context.Context, queue string) (*sqs.GetQueueUrlOutput, error) {
	return s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
}

func (s *sqsPoller) DeleteMessage(ctx context.Context, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	input.QueueUrl = s.queueURL
	return s.client.DeleteMessage(ctx, input)
}

func (s *sqsPoller) GetQueueAttrs(ctx context.Context, input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	input.QueueUrl = s.queueURL
	return s.client.GetQueueAttributes(ctx, input)
}
