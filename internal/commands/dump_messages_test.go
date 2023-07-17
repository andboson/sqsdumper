package commands

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	mock_aws "andboson/sqsdumper/internal/mocks/mock_sqs"
	"andboson/sqsdumper/internal/wrappers/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/ptr"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var log = zerolog.New(os.Stderr).With().Logger()

func TestSQSDumper_ProcessMessagesFull(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	rawMessage := json.RawMessage(`{"foo":"bar"}`)
	msg := types.Message{
		Body: getBody(t, aws.EventMessage{
			Message: &rawMessage,
		}),
		MessageId: nil,
	}

	poller := mock_aws.NewMockSQSPoller(ctrl)
	poller.EXPECT().DeleteMessage(gomock.Any(), gomock.Any())
	poller.EXPECT().GetQueueURL().Return(ptr.String("url"))

	params := SQSDumperParams{
		Logger:        log,
		DeleteMessage: true,
		RawMessage:    true,
		JsonPath:      "",
	}

	dumper := NewSQSDumper(params)
	err := dumper.ProcessMessages(ctx)(poller, msg)
	assert.NoError(t, err)
}

func TestSQSDumper_ProcessMessagesSimple(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	msg := types.Message{
		Body:      ptr.String(`{"foo":"bar"}`),
		MessageId: nil,
	}

	poller := mock_aws.NewMockSQSPoller(ctrl)
	params := SQSDumperParams{
		Logger:        log,
		DeleteMessage: false,
		RawMessage:    false,
		JsonPath:      "",
	}
	dumper := NewSQSDumper(params)
	err := dumper.ProcessMessages(ctx)(poller, msg)
	assert.NoError(t, err)
}

func getBody(t *testing.T, message aws.EventMessage) *string {
	b, err := json.Marshal(message)
	str := string(b)
	assert.NoError(t, err)

	return &str
}
