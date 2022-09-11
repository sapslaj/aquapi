package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sapslaj/aquapi/internal/maintenance"
)

func handler(ctx context.Context) error {
	return maintenance.SyncImagesBucketToTable(ctx)
}

func main() {
	lambda.Start(handler)
}
