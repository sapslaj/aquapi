package graph

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/aquapics"
)

type ImagesSort interface {
	GetNext() (types.Object, error)
}

type ImagesSortRandom struct{}

func (is *ImagesSortRandom) GetNext() (types.Object, error) {
	return aquapics.GetRandomFromS3()
}

type ImagesSortKey struct {
	cache      []types.Object
	cacheIndex int
	prefix     string
	startAfter string
}

func (is *ImagesSortKey) warmCache() error {
	if is.cacheIndex < len(is.cache) {
		return nil
	}
	objects, err := aquapics.ListFromS3(is.prefix, is.startAfter, 10)
	if err != nil {
		return err
	}
	if len(objects) == 0 {
		return errors.New("result list empty")
	}
	is.startAfter = *objects[len(objects)-1].Key
	is.cacheIndex = 0
	is.cache = objects
	return nil
}

func (is *ImagesSortKey) GetNext() (types.Object, error) {
	err := is.warmCache()
	if err != nil {
		return types.Object{}, err
	}
	is.cacheIndex += 1
	return is.cache[is.cacheIndex-1], nil
}

func NewImagesSort(sort string, afterKey string) (ImagesSort, error) {
	switch sort {
	case "random":
		return &ImagesSortRandom{}, nil
	case "key":
		return &ImagesSortKey{
			cacheIndex: 0,
			prefix:     "",
			startAfter: afterKey,
		}, nil
	default:
		return nil, fmt.Errorf("sort method '%s' is not supported", sort)
	}
}
