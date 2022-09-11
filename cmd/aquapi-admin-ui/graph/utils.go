package graph

import (
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/model"
	"github.com/sapslaj/aquapi/internal/service"
)

func imageToModel(image *service.Image) (*model.Image, error) {
	url, err := image.GetUrl()
	if err != nil {
		return nil, err
	}
	return &model.Image{
		ID:   image.ID,
		URL:  url,
		Tags: image.Tags,
	}, nil
}

func imageSliceContainsId(l []*model.Image, id string) bool {
	for _, m := range l {
		if m.ID == id {
			return true
		}
	}
	return false
}
