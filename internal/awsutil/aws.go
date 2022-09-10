package awsutil

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func DefaultAwsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.Region = os.Getenv("AWS_REGION")
		return nil
	})
}
