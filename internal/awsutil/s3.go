package awsutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DefaultS3ClientFromConfig(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func DefaultS3Client() (*s3.Client, error) {
	cfg, err := DefaultAwsConfig()
	return DefaultS3ClientFromConfig(cfg), err
}

func RegionalS3ClientFromConfig(cfg aws.Config, location string) *s3.Client {
	if location == "" {
		location = "us-east-1"
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = location
	})
}

func RegionalS3Client(location string) (*s3.Client, error) {
	cfg, err := DefaultAwsConfig()

	return RegionalS3ClientFromConfig(cfg, location), err
}

func GetS3ClientForBucketName(bucketName string) (*s3.Client, error) {
	awsBucketName := aws.String(bucketName)
	cfg, err := DefaultAwsConfig()
	if err != nil {
		return &s3.Client{}, err
	}
	defaultS3Client := DefaultS3ClientFromConfig(cfg)
	locationOutput, err := defaultS3Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: awsBucketName,
	})
	if err != nil {
		return defaultS3Client, err
	}
	location := string(locationOutput.LocationConstraint)
	return RegionalS3Client(location)
}
