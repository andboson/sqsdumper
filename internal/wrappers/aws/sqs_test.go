package aws

import (
	"context"
	"os"
	"testing"

	"andboson/sqsdumper/internal/mocks/mock_aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var log = zerolog.New(os.Stderr).With().Logger()

func TestSqsPoller_fetchQueueURL(t *testing.T) {
	ctx := context.Background()
	awsClient := NewAWSClient()

	cfg, err := awsClient.LoadDefaultConfig(ctx)
	assert.NoError(t, err)
	poller := &sqsPoller{
		client: sqs.NewFromConfig(cfg),
		logger: log,
	}
	assert.NotNil(t, poller)

	queueName := ""
	urlOutput, err := poller.fetchQueueURL(ctx, queueName)
	assert.Error(t, err)
	assert.Nil(t, urlOutput)
}

func TestSqsPoller_PollMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := gomock.NewController(t)

	t.Run("run success", func(t *testing.T) {
		sqsClient := mock_aws.NewMockSQSAPI(ctrl)

		sqsClient.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).
			Return(&sqs.GetQueueUrlOutput{}, nil)
		sqsClient.EXPECT().GetQueueAttributes(gomock.Any(), gomock.Any()).
			Return(&sqs.GetQueueAttributesOutput{
				Attributes: map[string]string{
					string(types.QueueAttributeNameApproximateNumberOfMessages): "1",
				},
			}, nil)

		msgID := "#1"
		var gotMessages int
		poller, err := NewSQSPoller(SQSParam{Client: sqsClient, Logger: log, QueueConfig: ConfigQueue{}})
		assert.NotNil(t, poller)
		// poll the messages
		sqsClient.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).
			Return(&sqs.ReceiveMessageOutput{
				Messages: []types.Message{{
					MessageId: &msgID,
				}},
			}, nil).Times(1)

		sqsClient.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ *sqs.ReceiveMessageInput, _ ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
				cancel()

				return nil, errors.New("some error")
			})

		err = poller.PollMessages(ctx, func(poller SQSPoller, msg types.Message) error {
			gotMessages++
			return nil
		})
		assert.NoError(t, err)

		<-ctx.Done()
		assert.Equal(t, 1, gotMessages)
	})
}

func TestSqsPoller_PollMessagesErrorProcessing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := gomock.NewController(t)

	t.Run("run success", func(t *testing.T) {
		sqsClient := mock_aws.NewMockSQSAPI(ctrl)

		sqsClient.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).
			Return(&sqs.GetQueueUrlOutput{}, nil)
		sqsClient.EXPECT().GetQueueAttributes(gomock.Any(), gomock.Any()).
			Return(&sqs.GetQueueAttributesOutput{
				Attributes: map[string]string{
					string(types.QueueAttributeNameApproximateNumberOfMessages): "1",
				},
			}, nil)

		msgID := "#1"
		var gotMessages int
		poller, err := NewSQSPoller(SQSParam{Client: sqsClient, Logger: log, QueueConfig: ConfigQueue{}})
		assert.NotNil(t, poller)

		// poll the messages
		sqsClient.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).
			Return(&sqs.ReceiveMessageOutput{
				Messages: []types.Message{{
					MessageId: &msgID,
				}},
			}, nil).Times(1)

		sqsClient.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ *sqs.ReceiveMessageInput, _ ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
				cancel()

				return &sqs.ReceiveMessageOutput{
					Messages: []types.Message{},
				}, nil
			})

		err = poller.PollMessages(ctx, func(poller SQSPoller, msg types.Message) error {
			gotMessages++
			return errors.New("some error")
		})
		assert.NoError(t, err)

		<-ctx.Done()
		assert.Equal(t, 1, gotMessages)
	})
}
