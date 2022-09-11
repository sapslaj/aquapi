package maintenance

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/awsutil"
	"github.com/sapslaj/aquapi/internal/config"
	"github.com/sapslaj/aquapi/internal/service"
)

func SyncImageObjectToTable(ctx context.Context, object s3types.Object) {
	image, err := service.NewImageFromS3Object(object)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	_, err = image.GetTags()
	if err != nil {
		log.Print(err)
		panic(err)
	}
	err = image.Update()
	if err != nil {
		log.Print(err)
		panic(err)
	}
	log.Print(image.ID)
}

func SyncImagesBucketToTable(ctx context.Context) error {
	s3BucketClient, err := awsutil.GetS3ClientForBucketName(config.ImagesBucketName())
	if err != nil {
		return err
	}
	imagesBucketPaginator := s3.NewListObjectsV2Paginator(s3BucketClient, &s3.ListObjectsV2Input{
		Bucket: aws.String(config.ImagesBucketName()),
	})
	for imagesBucketPaginator.HasMorePages() {
		var wg sync.WaitGroup
		output, err := imagesBucketPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, object := range output.Contents {
			wg.Add(1)
			go func(object s3types.Object) {
				defer wg.Done()
				SyncImageObjectToTable(ctx, object)
			}(object)
		}
		wg.Wait()
	}
	return nil
}
