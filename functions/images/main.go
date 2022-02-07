package main

import (
	"context"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sapslaj/aquapi/internal/api"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count := 1
	for key, value := range request.QueryStringParameters {
		if key == "count" {
			c, err := strconv.Atoi(value)
			if err != nil {
				return api.ResponseError(400, "`count` is unparseable", err.Error()), nil
			}
			if c >= 100 {
				return api.ResponseError(400, "`count` is too large", "only up to 100 images can be returned at a time."), nil
			}
			count = c
		}
	}
	images := make([]*api.Image, count)
	for c := 0; c < count; c++ {
		object, err := aquapics.GetRandomFromS3()
		if err != nil {
			return api.ResponseError(503, "An internal error occurred", err.Error()), err
		}
		images[c] = &api.Image{
			ID:  *object.Key,
			URL: utils.S3ObjectToUrl(object),
		}
	}

	return api.ResponseSuccess(images), nil
}

func main() {
	lambda.Start(handler)
}
