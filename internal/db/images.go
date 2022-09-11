package db

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/ptr"
	"github.com/sapslaj/aquapi/internal/awsutil"
	"github.com/sapslaj/aquapi/internal/config"
)

const idCharacters = "0123456789abcdef"

type Image struct {
	ID   string
	Tags []string
}

func randInt(max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64())
}

func randidCharacter() string {
	return string(idCharacters[randInt(len(idCharacters))])
}

func strSliceContains(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

func NewImageFromDDBItem(item map[string]dynamodbtypes.AttributeValue) *Image {
	var image *Image
	attributevalue.UnmarshalMap(item, &image)
	return image
}

func RandomImage(allowTags []*string, omitTags []*string) (*Image, error) {
	dynamodbClient, err := awsutil.DefaultDynamoDBClient()
	if err != nil {
		return nil, err
	}
	filterExpression := "begins_with(#id, :prefix)"
	expressionAttributeNames := map[string]string{
		"#id": "id",
	}
	expressionAttributeValues := map[string]dynamodbtypes.AttributeValue{}
	if allowTags != nil || omitTags != nil {
		expressionAttributeNames["#tags"] = "tags"
	}
	for tagIdx, tag := range ptr.ToStringSlice(omitTags) {
		attrValue := fmt.Sprintf(":omittag%d", tagIdx)
		filterExpression += fmt.Sprintf(" AND NOT contains(#tags, %s)", attrValue)
		expressionAttributeValues[attrValue] = &dynamodbtypes.AttributeValueMemberS{Value: tag}
	}
	if allowTags != nil {
		filterExpression += "AND ("
		for tagIdx, tag := range allowTags {
			attrValue := fmt.Sprintf(":onlytag%d", tagIdx)
			if tagIdx > 0 {
				filterExpression += " OR "
			}
			filterExpression += fmt.Sprintf("contains(#tags, %s)", attrValue)
			expressionAttributeValues[attrValue] = &dynamodbtypes.AttributeValueMemberS{Value: *tag}
		}
		filterExpression += ")"
	}
	for {
		prefix := randidCharacter()
		expressionAttributeValues[":prefix"] = &dynamodbtypes.AttributeValueMemberS{Value: prefix}
		scanPaginator := dynamodb.NewScanPaginator(dynamodbClient, &dynamodb.ScanInput{
			TableName:                 aws.String(config.ImagesDynamoDBTable()),
			FilterExpression:          aws.String(filterExpression),
			ExpressionAttributeNames:  expressionAttributeNames,
			ExpressionAttributeValues: expressionAttributeValues,
		})
		candidateItems := make([]map[string]dynamodbtypes.AttributeValue, 0)
		for scanPaginator.HasMorePages() {
			scanOutput, err := scanPaginator.NextPage(context.TODO())
			if err != nil {
				return nil, err
			}
			candidateItems = append(candidateItems, scanOutput.Items...)
		}
		if len(candidateItems) == 0 {
			continue
		}
		for attempt := 1; attempt < len(candidateItems); attempt++ {
			item := candidateItems[randInt(len(candidateItems))]
			image := NewImageFromDDBItem(item)
			acceptable := true
			if allowTags != nil {
				acceptable = false
				for _, tag := range image.Tags {
					if !acceptable && strSliceContains(ptr.ToStringSlice(allowTags), tag) {
						acceptable = true
					}
				}
			}
			if acceptable {
				return image, nil
			}
		}
	}
}

func (i *Image) CreateOrUpdate() (*dynamodb.PutItemOutput, error) {
	dynamodbClient, err := awsutil.DefaultDynamoDBClient()
	if err != nil {
		return nil, err
	}
	item := make(map[string]dynamodbtypes.AttributeValue)
	item["id"] = &dynamodbtypes.AttributeValueMemberS{Value: i.ID}
	if len(i.Tags) > 0 {
		item["tags"] = &dynamodbtypes.AttributeValueMemberSS{Value: i.Tags}
	}
	return dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(config.ImagesDynamoDBTable()),
		Item:      item,
	})
}
