package aquapics

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const keyCharacters = "0123456789abcdef"

var ignoredKeys = []string{"error.html", "index.html", "favicon.ico"}

func randint(i int) int {
	rand.Seed(time.Now().UnixMicro())
	return rand.Intn(i)
}

func isBadKey(key string) bool {
	for _, v := range ignoredKeys {
		if key == v {
			return true
		}
	}
	return false
}

func getRandomFromContents(objects []s3types.Object) s3types.Object {
	object := s3types.Object{Key: aws.String(ignoredKeys[0])}
	for isBadKey(*object.Key) {
		rand.Seed(time.Now().UnixMicro())
		index := rand.Intn(len(objects))
		object = objects[index]
	}
	return object
}

// GetRandomFromS3 returns a random object from the S3 bucket
func GetRandomFromS3() (s3types.Object, error) {
	// This implementation is shit, but it works since there's currently so few
	// images in S3. Eventually I'll come up with a better algo.
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.Region = os.Getenv("AWS_REGION")
		return nil
	})
	if err != nil {
		return s3types.Object{}, err
	}

	bucketName := aws.String(os.Getenv("AQUAPI_IMAGES_BUCKET"))

	defaultS3Client := s3.NewFromConfig(cfg)
	locationOutput, err := defaultS3Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: bucketName,
	})
	if err != nil {
		return s3types.Object{}, err
	}
	location := string(locationOutput.LocationConstraint)
	if location == "" {
		location = "us-east-1"
	}

	regionalS3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = location
	})
	randomPrefix := string(keyCharacters[randint(len(keyCharacters))])
	output, err := regionalS3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: aws.String(randomPrefix),
	})
	if err != nil {
		return s3types.Object{}, err
	}

	return getRandomFromContents(output.Contents), nil
}
