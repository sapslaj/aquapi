package aquapics

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const tagsTagKey = "AquaPITags"
const tagSeperator = ":"

func GetTags(object s3types.Object) ([]string, error) {
	s3BucketClient, err := getS3ClientForBucketName(imagesBucketName)
	tags := []string{}
	if err != nil {
		return tags, err
	}
	output, err := s3BucketClient.GetObjectTagging(context.TODO(), &s3.GetObjectTaggingInput{
		Bucket: aws.String(imagesBucketName),
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
	s3BucketClient, err := getS3ClientForBucketName(imagesBucketName)
	if err != nil {
		return err
	}
	value := strings.Join(tags, tagSeperator)
	_, err = s3BucketClient.PutObjectTagging(context.TODO(), &s3.PutObjectTaggingInput{
		Bucket: aws.String(imagesBucketName),
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

func AddTag(object s3types.Object, tag string) error {
	existingTags, err := GetTags(object)
	if err != nil {
		return err
	}
	for _, existingTag := range existingTags {
		if existingTag == tag {
			return nil
		}
	}
	tags := append(existingTags, tag)
	return SetTags(object, tags)
}

func RemoveTag(object s3types.Object, tag string) error {
	existingTags, err := GetTags(object)
	if err != nil {
		return err
	}
	for index, existingTag := range existingTags {
		if existingTag == tag {
			tags := append(existingTags[:index], existingTags[index+1:]...)
			return SetTags(object, tags)
		}
	}
	return nil
}
