package aquapics

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ListFromS3(prefix, startAfter string, limit int) ([]s3types.Object, error) {
	s3BucketClient, err := getS3ClientForBucketName(imagesBucketName)
	if err != nil {
		return nil, err
	}
	output, err := s3BucketClient.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:     aws.String(imagesBucketName),
		Prefix:     aws.String(prefix),
		StartAfter: aws.String(startAfter),
		MaxKeys:    int32(limit),
	})
	if err != nil {
		return nil, err
	}
	return output.Contents, nil
}
