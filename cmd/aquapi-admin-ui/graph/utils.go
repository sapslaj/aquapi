package graph

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/ptr"
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/model"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/utils"
)

func imageObjectAndTagsToModel(object types.Object, tags []string) *model.Image {
	return &model.Image{
		ID:   *object.Key,
		URL:  utils.S3ObjectToUrl(object),
		Tags: tags,
	}
}

func imageObjectToModel(object types.Object) (*model.Image, error) {
	tags, err := aquapics.GetTags(object)
	if err != nil {
		return nil, err
	}
	return imageObjectAndTagsToModel(object, tags), nil
}

func imageTagAction(input *model.ImageTagsInput, action func(types.Object) error) (*model.Image, error) {
	object, err := aquapics.GetImageObject(input.ImageID)
	if err != nil {
		return nil, err
	}
	err = action(object)
	if err != nil {
		return nil, err
	}
	return imageObjectToModel(object)
}

func imageSliceContainsId(l []*model.Image, id string) bool {
	for _, m := range l {
		if m.ID == id {
			return true
		}
	}
	return false
}

func strSliceContains(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

func applicableImage(tags []string, allowTags []*string, omitTags []*string, onlyTags []*string) bool {
	if onlyTags == nil {
		if allowTags == nil && omitTags == nil {
			return true
		}
		acceptable := false
		if allowTags == nil {
			acceptable = true
		}
		for _, tag := range tags {
			if !acceptable && strSliceContains(ptr.ToStringSlice(allowTags), tag) {
				acceptable = true
			}
			if acceptable && strSliceContains(ptr.ToStringSlice(omitTags), tag) {
				acceptable = false
			}
		}
		return acceptable
	} else {
		for _, onlyTag := range onlyTags {
			if !strSliceContains(tags, *onlyTag) {
				return false
			}
		}
		for _, tag := range tags {
			if !strSliceContains(ptr.ToStringSlice(onlyTags), tag) {
				return false
			}
		}
		return true
	}
}
