package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/ptr"
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/generated"
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/model"
	"github.com/sapslaj/aquapi/internal/aquapics"
)

func (r *mutationResolver) AddTags(ctx context.Context, input *model.ImageTagsInput) (*model.Image, error) {
	return imageTagAction(input, func(object types.Object) error {
		for _, tag := range input.Tags {
			if err := aquapics.AddTag(object, tag); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *mutationResolver) RemoveTags(ctx context.Context, input *model.ImageTagsInput) (*model.Image, error) {
	return imageTagAction(input, func(object types.Object) error {
		for _, tag := range input.Tags {
			if err := aquapics.RemoveTag(object, tag); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *mutationResolver) SetTags(ctx context.Context, input *model.ImageTagsInput) (*model.Image, error) {
	return imageTagAction(input, func(object types.Object) error {
		return aquapics.SetTags(object, input.Tags)
	})
}

func (r *queryResolver) Image(ctx context.Context, id string) (*model.Image, error) {
	object, err := aquapics.GetImageObject(id)
	if err != nil {
		return nil, err
	}
	return imageObjectToModel(object)
}

func (r *queryResolver) Images(ctx context.Context, sort *string, limit *int, afterKey *string, allowTags []*string, omitTags []*string, onlyTags []*string) ([]*model.Image, error) {
	if sort == nil {
		sort = ptr.String("random")
	}
	if afterKey == nil {
		afterKey = ptr.String("")
	}
	imagesSort, err := NewImagesSort(*sort, *afterKey)
	if err != nil {
		return nil, err
	}
	defaultLimit := 10
	if limit == nil {
		limit = &defaultLimit
	}
	if *limit >= 100 {
		return nil, fmt.Errorf("lower limit pls")
	}
	models := []*model.Image{}
	for len(models) < *limit {
		object, err := imagesSort.GetNext()
		if err != nil {
			return nil, err
		}
		if imageSliceContainsId(models, *object.Key) {
			continue
		}
		tags, err := aquapics.GetTags(object)
		if err != nil {
			return nil, err
		}
		if applicableImage(tags, allowTags, omitTags, onlyTags) {
			models = append(models, imageObjectAndTagsToModel(object, tags))
		}
	}
	return models, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
