package main

import (
	"context"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/smithy-go/ptr"
	"github.com/sapslaj/aquapi/internal/api"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/service"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count := 1
	nsfw := "none"
	imageService := service.NewImagesService()
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
	var allowTagsInput []*string
	var omitTagsInput []*string
	omitTags := []string{aquapics.HIDDEN, aquapics.MEME, aquapics.COLLAGE}
	if nsfw == "none" {
		omitTags = append(omitTags, aquapics.NSFW, aquapics.ECCHI, aquapics.HENTAI)
	} else if nsfw == "only" {
		allowTags := []string{aquapics.ECCHI, aquapics.HENTAI}
		allowTagsInput = ptr.StringSlice(allowTags)
	}
	omitTagsInput = ptr.StringSlice(omitTags)
	for len(images) < count {
		image, err := imageService.GetRandomImageFilterTags(allowTagsInput, omitTagsInput)
		if err != nil {
			return api.ResponseError(503, "An internal error occurred", err.Error()), err
		}
		imageDto, err := api.NewImageFromImagesServiceImage(image)
		if err != nil {
			return api.ResponseError(503, "An internal error occurred", err.Error()), err
		}
		if err != nil {
			return api.ResponseError(503, "An internal error occurred", err.Error()), err
		}
		images = append(images, imageDto)
	}

	return api.ResponseSuccess(images), nil
}

func main() {
	lambda.Start(handler)
}
