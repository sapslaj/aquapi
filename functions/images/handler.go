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

func applicableImage(tags []string, nsfw string) bool {
	hasNsfwTag := false
	for _, tag := range tags {
		switch tag {
		case "hidden", "meme", "collage":
			return false
		case "nsfw", "ecchi", "hentai":
			hasNsfwTag = true
		}
	}
	if hasNsfwTag && nsfw == "none" {
		return false
	} else if hasNsfwTag && nsfw == "allowed" {
		return true
	} else if !hasNsfwTag && nsfw == "only" {
		return false
	}
	return true
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count := 1
	nsfw := "none"
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
		if key == "nsfw" {
			valid := false
			for _, v := range []string{"none", "allowed", "only"} {
				if v == value {
					valid = true
					break
				}
			}
			if !valid {
				return api.ResponseError(400, "`nsfw` is an invalid value", "`nsfw` must equal 'none', 'allowed', or 'only'."), nil
			}
			nsfw = value
		}
	}
	images := []*api.Image{}
	for len(images) < count {
		object, err := aquapics.GetRandomFromS3()
		if err != nil {
			return api.ResponseError(503, "An internal error occurred", err.Error()), err
		}
		tags, err := aquapics.GetTags(object)
		if err != nil {
			return api.ResponseError(503, "An internal error occurred", err.Error()), err
		}
		if applicableImage(tags, nsfw) {
			images = append(images, &api.Image{
				ID:   *object.Key,
				URL:  utils.S3ObjectToUrl(object),
				Tags: tags,
			})
		}
	}

	return api.ResponseSuccess(images), nil
}

func main() {
	lambda.Start(handler)
}
