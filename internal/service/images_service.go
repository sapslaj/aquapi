package service

import (
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/db"
	"github.com/sapslaj/aquapi/internal/utils"
)

type ImagesService struct {
}

type Image struct {
	ID   string
	Tags []string

	dbo      *db.Image
	s3Object *s3types.Object
}

func NewImagesService() *ImagesService {
	return &ImagesService{}
}

func (is *ImagesService) FindById(id string) (*Image, error) {
	dbo, err := db.GetImageById(id)
	if err != nil {
		return nil, err
	}
	return NewImageFromImageDbo(dbo)
}

func (is *ImagesService) GetRandomImageFilterTags(allowTags []*string, omitTags []*string) (*Image, error) {
	dbo, err := db.RandomImage(allowTags, omitTags)
	if err != nil {
		return nil, err
	}
	return NewImageFromImageDbo(dbo)
}

func NewImageFromImageDbo(dbo *db.Image) (*Image, error) {
	return &Image{
		ID:   dbo.ID,
		Tags: dbo.Tags,
		dbo:  dbo,
	}, nil
}

func NewImageFromS3Object(s3Object s3types.Object) (*Image, error) {
	tags, err := aquapics.GetTags(s3Object)
	if err != nil {
		return nil, err
	}
	return &Image{
		ID:       *s3Object.Key,
		Tags:     tags,
		s3Object: &s3Object,
	}, nil
}

func (i *Image) Update() error {
	object, err := i.GetS3Object()
	if err != nil {
		return err
	}
	dbo, err := i.GetDBO()
	if err != nil {
		return err
	}
	err = aquapics.SetTags(*object, i.Tags)
	if err != nil {
		return err
	}
	dbo.Tags = i.Tags
	_, err = dbo.CreateOrUpdate()
	return err
}

func (i *Image) GetTags() ([]string, error) {
	s3Object, err := i.GetS3Object()
	if err != nil {
		return i.Tags, err
	}
	tags, err := aquapics.GetTags(*s3Object)
	if err != nil {
		return i.Tags, err
	}
	i.Tags = tags
	return i.Tags, nil
}

func (i *Image) SetTags(tags []string) error {
	i.Tags = tags
	return nil
}

func (i *Image) AddTag(tag string) error {
	for _, existingTag := range i.Tags {
		if existingTag == tag {
			return nil
		}
	}
	i.Tags = append(i.Tags, tag)
	return nil
}

func (i *Image) RemoveTag(tag string) error {
	for index, existingTag := range i.Tags {
		if existingTag == tag {
			tags := append(i.Tags[:index], i.Tags[index+1:]...)
			i.Tags = tags
		}
	}
	return nil
}

func (i *Image) GetUrl() (string, error) {
	return utils.ImagesIDToUrl(i.ID), nil
}

func (i *Image) GetDBO() (*db.Image, error) {
	if i.dbo == nil {
		dbo, err := db.GetImageById(i.ID)
		if err != nil {
			return nil, err
		}
		i.dbo = dbo
	}
	return i.dbo, nil
}

func (i *Image) GetS3Object() (*s3types.Object, error) {
	if i.s3Object == nil {
		s3Object, err := aquapics.GetImageObject(i.ID)
		if err != nil {
			return nil, err
		}
		i.s3Object = &s3Object
	}
	return i.s3Object, nil
}
