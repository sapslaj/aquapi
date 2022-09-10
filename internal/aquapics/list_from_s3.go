package aquapics

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/awsutil"
)

func ListFromS3(prefix, startAfter string, limit int) ([]s3types.Object, error) {
	s3BucketClient, err := awsutil.GetS3ClientForBucketName(ImagesBucketName)
	if err != nil {
		return nil, err
	}
	output, err := s3BucketClient.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:     aws.String(ImagesBucketName),
		Prefix:     aws.String(prefix),
		StartAfter: aws.String(startAfter),
		MaxKeys:    int32(limit),
	})
	if err != nil {
		return nil, err
	}
	return output.Contents, nil
}
