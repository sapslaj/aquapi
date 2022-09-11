package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/generated"
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/model"
	"github.com/sapslaj/aquapi/internal/service"
)

func (r *mutationResolver) AddTags(ctx context.Context, input *model.ImageTagsInput) (*model.Image, error) {
	imagesService := service.NewImagesService()
	image, err := imagesService.FindById(input.ImageID)
	if err != nil {
		return nil, err
	}
	for _, tag := range input.Tags {
		if err := image.AddTag(tag); err != nil {
			return nil, err
		}
	}
	image.Update()
	return imageToModel(image)
}

func (r *mutationResolver) RemoveTags(ctx context.Context, input *model.ImageTagsInput) (*model.Image, error) {
	imagesService := service.NewImagesService()
	image, err := imagesService.FindById(input.ImageID)
	if err != nil {
		return nil, err
	}
	for _, tag := range input.Tags {
		if err := image.RemoveTag(tag); err != nil {
			return nil, err
		}
	}
	if err := image.Update(); err != nil {
		return nil, err
	}
	return imageToModel(image)
}

func (r *mutationResolver) SetTags(ctx context.Context, input *model.ImageTagsInput) (*model.Image, error) {
	imagesService := service.NewImagesService()
	image, err := imagesService.FindById(input.ImageID)
	if err != nil {
		return nil, err
	}
	if err := image.SetTags(input.Tags); err != nil {
		return nil, err
	}
	if err := image.Update(); err != nil {
		return nil, err
	}
	return imageToModel(image)
}

func (r *queryResolver) Image(ctx context.Context, id string) (*model.Image, error) {
	imagesService := service.NewImagesService()
	image, err := imagesService.FindById(id)
	if err != nil {
		return nil, err
	}
	return imageToModel(image)
}

func (r *queryResolver) Images(ctx context.Context, limit *int, allowTags []*string, omitTags []*string) ([]*model.Image, error) {
	imagesService := service.NewImagesService()
	defaultLimit := 10
	if limit == nil {
		limit = &defaultLimit
	}
	if *limit >= 100 {
		return nil, fmt.Errorf("lower limit pls")
	}
	models := []*model.Image{}
	for len(models) < *limit {
		image, err := imagesService.GetRandomImageFilterTags(allowTags, omitTags)
		if err != nil {
			return nil, err
		}
		if imageSliceContainsId(models, image.ID) {
			continue
		}
		model, err := imageToModel(image)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
