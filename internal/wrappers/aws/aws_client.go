package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
)

// Client represents AWS client
//
//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_aws/mock_$GOFILE
type Client interface {
	LoadDefaultConfig(ctx context.Context) (aws.Config, error)
	LoadLocalConfig(ctx context.Context) (aws.Config, error)
}

// NewAWSClient returns an instance of service
func NewAWSClient() Client {
	return &awsClient{}
}

type awsClient struct {
	config aws.Config
}

// LoadDefaultConfig loads and returns default config for AWS from ENV or the  ./aws config location
func (a *awsClient) LoadDefaultConfig(ctx context.Context) (aws.Config, error) {
	if os.Getenv("localstack") != "" {
		fmt.Printf("\n== LOADING LOCAL CONFIG ==\n")
		return a.LoadLocalConfig(ctx)
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), awsMaxAttempts)
	}))
	if err != nil {
		return cfg, errors.Wrap(err, "configuration error")
	}

	a.config = cfg

	return cfg, nil
}

// LoadLocalConfig loads config with a local endpoint, for a testing purpose
func (a *awsClient) LoadLocalConfig(_ context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               "http://localstack:4566",
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			})),
	)
	if err != nil {
		return cfg, errors.Wrap(err, "configuration error")
	}

	a.config = cfg

	return cfg, nil
}
