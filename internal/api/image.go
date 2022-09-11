package api

import "github.com/sapslaj/aquapi/internal/service"

type Image struct {
	ID   string   `jsonapi:"primary,image"`
	URL  string   `jsonapi:"attr,url"`
	Tags []string `jsonapi:"attr,tags"`
}

func NewImageFromImagesServiceImage(image *service.Image) (*Image, error) {
	url, err := image.GetUrl()
	if err != nil {
		return nil, err
	}
	return &Image{
		ID:   image.ID,
		URL:  url,
		Tags: image.Tags,
	}, nil
}
