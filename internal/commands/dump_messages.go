package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"andboson/sqsdumper/internal/wrappers/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sinhashubham95/jsonic"
)

type SQSDumperParams struct {
	Logger        zerolog.Logger
	DeleteMessage bool
	RawMessage    bool
	JsonPath      string
}

// SQSDumper is a command to print a message content
type SQSDumper struct {
	logger        zerolog.Logger
	deleteMessage bool
	rawMessage    bool
	jsonPath      string
}

// NewSQSDumper returns a new instance
func NewSQSDumper(p SQSDumperParams) SQSDumper {
	return SQSDumper{
		logger:        p.Logger,
		deleteMessage: p.DeleteMessage,
		rawMessage:    p.RawMessage,
		jsonPath:      p.JsonPath,
	}
}

// ProcessMessages returns aws.MessageHandler type func which process the incoming message
func (p *SQSDumper) ProcessMessages(ctx context.Context) aws.MessageHandler {
	p.logger.Info().Msg("started processing")
	return func(sqsPoller aws.SQSPoller, msg types.Message) error {
		// process the message
		if err := p.processMessage(ctx, msg); err != nil {
			p.logger.Err(err).Msg("error process the message")
			return errors.Wrapf(err, "error processing the message")
		}

		if !p.deleteMessage {
			return nil
		}

		// delete the message
		if _, err := sqsPoller.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      sqsPoller.GetQueueURL(),
			ReceiptHandle: msg.ReceiptHandle,
		}); err != nil {
			p.logger.Err(err).Msg("error deleting the message")
			return errors.Wrapf(err, "error deleting the message")
		}

		return nil
	}
}

func (p *SQSDumper) processMessage(_ context.Context, msg types.Message) error {
	eventMessage, err := aws.ParseEventMessage(*msg.Body)
	if err != nil {
		return errors.Wrap(err, "error parsing the incoming message")
	}

	if p.rawMessage || eventMessage.Message == nil {
		fmt.Println(*msg.Body)
		return nil
	}

	stringed := string(*eventMessage.Message)

	if p.jsonPath != "" {
		return p.printByPath(stringed)
	}

	stringed = strings.ReplaceAll(stringed, `\"`, `"`)
	if len(stringed) >= 2 {
		fmt.Println(stringed[1 : len(stringed)-1])
	} else {
		fmt.Println(stringed)
	}

	return nil
}

func (p *SQSDumper) printByPath(msg string) error {
	str, err := strconv.Unquote(msg)
	if err != nil {
		// just keep the message
		str = msg
	}

	j, err := jsonic.New([]byte(str))
	if err != nil {
		return err
	}

	data, err := j.Get(p.jsonPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", data)

	return nil
}
