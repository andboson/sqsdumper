package main

import (
	"fmt"
	"os"

	"andboson/sqsdumper/internal/commands"
	"andboson/sqsdumper/internal/wrappers/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rs/zerolog"
	cli "github.com/urfave/cli/v2"
)

// Version holds the application version
var Version string

func main() {
	var (
		stopOnTotal   *bool
		deleteMessage bool
		queueName     string
		jsonPath      string
	)

	app := &cli.App{
		Name:                 "sqsdumper",
		Version:              Version,
		Usage:                "sqsdumper -s src_queue",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "stopOnTotal",
				Usage:       "stop when all messages processed",
				Destination: stopOnTotal,
				DefaultText: "true",
			},
			&cli.BoolFlag{
				Name:        "deleteMessage",
				Usage:       "delete received messages",
				Destination: &deleteMessage,
			},
			&cli.StringFlag{
				Name:        "queueName",
				Aliases:     []string{"s"},
				Usage:       "the source queue",
				Destination: &queueName,
			},
			&cli.StringFlag{
				Name:        "jsonPath",
				Aliases:     []string{"jp"},
				Usage:       "json path, like x.y, see https://github.com/sinhashubham95/jsonic for more",
				Destination: &jsonPath,
				DefaultText: ".",
			},
		},
		Action: func(ctx *cli.Context) error {
			l := zerolog.New(os.Stderr).With().Timestamp().Logger()
			commander := commands.NewSQSDumper(l, deleteMessage, jsonPath)

			// Init AWS
			client := aws.NewAWSClient()
			cfg, err := client.LoadDefaultConfig(ctx.Context)
			if err != nil {
				l.Err(err).Msg("can't load the AWS config")
				return err
			}

			stop := true
			if stopOnTotal != nil {
				stop = *stopOnTotal
			}

			// Init BCQueue client and run poller
			poller, err := aws.NewSQSPoller(
				aws.SQSParam{
					Client: sqs.NewFromConfig(cfg),
					Logger: l,
					QueueConfig: aws.ConfigQueue{
						QueueName:               queueName,
						MaxMessagesPerRetrieval: 2,
						WaitTimeSeconds:         2,
					},
					StopOnTotal: stop,
					CounterChan: nil,
				},
			)
			if err != nil {
				l.Err(err).Msg("error creating SQS poller")
				return err
			}
			total := poller.GetTotal()

			defer func() {
				l.Log().Msgf(" === total processed: %d", total)
			}()

			return poller.PollMessages(ctx.Context, commander.ProcessMessages(ctx.Context))
		},
		Before: func(context *cli.Context) error {
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
