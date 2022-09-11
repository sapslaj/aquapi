package config

import "os"

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("key does not exist: " + key)
	}
	return value
}

func ImagesBucketName() string {
	return getEnv("AQUAPI_IMAGES_BUCKET")
}

func ImagesDynamoDBTable() string {
	return "aquapi-images-" + getEnv("AQUAPI_STAGE")
}
