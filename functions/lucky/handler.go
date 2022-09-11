package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/smithy-go/ptr"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/service"
)

// Response maps to APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

func handler(ctx context.Context) (Response, error) {
	imagesService := service.NewImagesService()
	omitTags := ptr.StringSlice(aquapics.TAGS)
	image, err := imagesService.GetRandomImageFilterTags(nil, omitTags)
	if err != nil {
		return Response{StatusCode: 503}, err
	}
	url, err := image.GetUrl()
	if err != nil {
		return Response{StatusCode: 503}, err
	}
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
