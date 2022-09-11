package maintenance

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/sapslaj/aquapi/internal/awsutil"
	"github.com/sapslaj/aquapi/internal/config"
	"github.com/sapslaj/aquapi/internal/db"
)

func SyncImagesBucketToTable(ctx context.Context) error {
	s3BucketClient, err := awsutil.GetS3ClientForBucketName(config.ImagesBucketName())
	if err != nil {
		return err
	}
	imagesBucketPaginator := s3.NewListObjectsV2Paginator(s3BucketClient, &s3.ListObjectsV2Input{
		Bucket: aws.String(config.ImagesBucketName()),
	})
	for imagesBucketPaginator.HasMorePages() {
		output, err := imagesBucketPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		for _, object := range output.Contents {
			wg.Add(1)
			go func(object s3types.Object) {
				defer wg.Done()
				if aquapics.IsBadKey(*object.Key) {
					return
				}
				tags, err := aquapics.GetTags(object)
				if err != nil {
					log.Print(err)
					panic(err)
				}
				image := db.Image{
					ID:   aws.ToString(object.Key),
					Tags: tags,
				}
				_, err = image.CreateOrUpdate()
				if err != nil {
					log.Print(err)
					panic(err)
				}
				log.Print(*object.Key)
			}(object)
		}
		wg.Wait()
	}
	return nil
}
