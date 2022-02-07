package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/utils"
)

// Respone maps to APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

func handler(ctx context.Context) (Response, error) {
	object, err := aquapics.GetRandomFromS3()
	if err != nil {
		return Response{StatusCode: 503}, err
	}
	url := utils.S3ObjectToUrl(object)
	body, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		return Response{StatusCode: 503}, err
	}
	return Response{
		StatusCode: 302,
		Body:       string(body),
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
