package aquapics

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/awsutil"
	"github.com/sapslaj/aquapi/internal/config"
)

const tagsTagKey = "AquaPITags"
const tagSeperator = ":"

// Tag definitions
const (
	HIDDEN  string = "hidden"
	NSFW    string = "nsfw"
	ECCHI   string = "ecchi"
	HENTAI  string = "hentai"
	MEME    string = "meme"
	COLLAGE string = "collage"
)

var TAGS = []string{
	HIDDEN,
	NSFW,
	ECCHI,
	HENTAI,
	MEME,
	COLLAGE,
}

func GetTags(object s3types.Object) ([]string, error) {
	s3BucketClient, err := awsutil.GetS3ClientForBucketName(config.ImagesBucketName())
	tags := []string{}
	if err != nil {
		return tags, err
	}
	output, err := s3BucketClient.GetObjectTagging(context.TODO(), &s3.GetObjectTaggingInput{
		Bucket: aws.String(config.ImagesBucketName()),
		Key:    object.Key,
	})
	if err != nil {
		return tags, err
	}
	for _, tag := range output.TagSet {
		if aws.ToString(tag.Key) == tagsTagKey {
			raw := aws.ToString(tag.Value)
			for _, tag := range strings.Split(raw, tagSeperator) {
				if tag == "" {
					continue
				}
				tags = append(tags, tag)
			}
			break
		}
	}
	return tags, nil
}

func SetTags(object s3types.Object, tags []string) error {
	s3BucketClient, err := awsutil.GetS3ClientForBucketName(config.ImagesBucketName())
	if err != nil {
		return err
	}
	value := strings.Join(tags, tagSeperator)
	_, err = s3BucketClient.PutObjectTagging(context.TODO(), &s3.PutObjectTaggingInput{
		Bucket: aws.String(config.ImagesBucketName()),
		Key:    object.Key,
		Tagging: &s3types.Tagging{
			TagSet: []s3types.Tag{
				{
					Key:   aws.String(tagsTagKey),
					Value: aws.String(value),
				},
			},
		},
	})
	return err
}
