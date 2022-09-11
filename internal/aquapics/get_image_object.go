package aquapics

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/awsutil"
	"github.com/sapslaj/aquapi/internal/config"
)

func GetImageObject(key string) (s3types.Object, error) {
	s3BucketClient, err := awsutil.GetS3ClientForBucketName(config.ImagesBucketName())
	if err != nil {
		return s3types.Object{}, err
	}
	output, err := s3BucketClient.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  aws.String(config.ImagesBucketName()),
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
