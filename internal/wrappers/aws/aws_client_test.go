package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAwsClient_LoadDefaultConfig(t *testing.T) {
	//from env
	const region = "eu-central-1"
	os.Setenv("AWS_REGION", region)

	client := NewAWSClient()
	cfg, err := client.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, cfg.Region, region)

	//from file
	const (
		region2    = "eu-central-2"
		configFile = "/tmp/aws_cfg.yaml"
	)
	os.Unsetenv("AWS_REGION") // without unsetting only already defined ENV-var will be used
	os.Setenv("AWS_CONFIG_FILE", configFile)
	err = createConfig(configFile, region2)
	assert.NoError(t, err)

	client = NewAWSClient()
	cfg, err = client.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, cfg.Region, region2)

	// can't load
	os.Setenv("AWS_ENABLE_ENDPOINT_DISCOVERY", "some bad value")
	client = NewAWSClient()
	_, err = client.LoadDefaultConfig(context.TODO())
	assert.NotEmpty(t, err)
	os.Unsetenv("AWS_ENABLE_ENDPOINT_DISCOVERY")
	os.Unsetenv("AWS_CONFIG_FILE")
}

func TestAwsClient_LoadLocalConfig(t *testing.T) {
	//from env
	const region = "eu-central-1"
	os.Setenv("AWS_REGION", region)

	client := NewAWSClient()
	cfg, err := client.LoadLocalConfig(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, cfg.Region, region)
}

func createConfig(path, region string) error {
	cont := fmt.Sprintf(`
[default]
region = %s
`, region)

	if err := ioutil.WriteFile(path, []byte(cont), 0777); err != nil {
		return errors.Wrap(err, "save error")
	}

	return nil
}
