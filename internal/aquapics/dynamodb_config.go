package aquapics

import "os"

var ImagesDynamoDBTable = "aquapi-images-" + os.Getenv("AQUAPI_STAGE")
