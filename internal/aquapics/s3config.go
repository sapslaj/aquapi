package aquapics

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var imagesBucketName = os.Getenv("AQUAPI_IMAGES_BUCKET")

func defaultAwsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.Region = os.Getenv("AWS_REGION")
		return nil
	})
}

func defaultS3ClientFromConfig(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func defaultS3Client() (*s3.Client, error) {
	cfg, err := defaultAwsConfig()
	return defaultS3ClientFromConfig(cfg), err
}

func regionalS3ClientFromConfig(cfg aws.Config, location string) *s3.Client {
	if location == "" {
		location = "us-east-1"
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = location
	})
}

func regionalS3Client(location string) (*s3.Client, error) {
	cfg, err := defaultAwsConfig()

	return regionalS3ClientFromConfig(cfg, location), err
}

func getS3ClientForBucketName(bucketName string) (*s3.Client, error) {
	awsBucketName := aws.String(bucketName)
	cfg, err := defaultAwsConfig()
	if err != nil {
		return &s3.Client{}, err
	}
	defaultS3Client := defaultS3ClientFromConfig(cfg)
	locationOutput, err := defaultS3Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: awsBucketName,
	})
	if err != nil {
		return defaultS3Client, err
	}
	location := string(locationOutput.LocationConstraint)
	return regionalS3Client(location)
}
