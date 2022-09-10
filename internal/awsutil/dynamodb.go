package awsutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DefaultDynamoDBClientFromConfig(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func DefaultDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := DefaultAwsConfig()
	return DefaultDynamoDBClientFromConfig(cfg), err
}
