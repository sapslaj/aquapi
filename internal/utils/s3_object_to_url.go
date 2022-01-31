package utils

import (
	"os"

	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func S3ObjectToUrl(object s3types.Object) string {
	imagesHost := os.Getenv("AQUAPI_IMAGES_HOST")
	return "https://" + imagesHost + "/" + *object.Key
}
