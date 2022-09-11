package utils

import (
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/config"
)

func ImagesIDToUrl(id string) string {
	return "https://" + config.ImagesHost() + "/" + id
}

func S3ObjectToUrl(object s3types.Object) string {
	return ImagesIDToUrl(*object.Key)
}
