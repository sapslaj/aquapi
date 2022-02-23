package aquapics

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func GetImageObject(key string) (s3types.Object, error) {
	s3BucketClient, err := getS3ClientForBucketName(imagesBucketName)
	if err != nil {
		return s3types.Object{}, err
	}
	output, err := s3BucketClient.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  aws.String(imagesBucketName),
		Prefix:  aws.String(key),
		MaxKeys: 1,
	})
	object := s3types.Object{}
	for _, o := range output.Contents {
		object = o
		break
	}
	return object, err
}
