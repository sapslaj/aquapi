package api

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/jsonapi"
)

// ResponseErrors is a factory function to generate a JSON API response with
// multiple ErrorObjects
func ResponseErrors(statusCode int, errors []*jsonapi.ErrorObject) events.APIGatewayProxyResponse {
	buf := new(bytes.Buffer)
	err := jsonapi.MarshalErrors(buf, errors)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503, Body: err.Error()}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       buf.String(),
	}
}

// ResponseError is a factory function to generate a JSON API response with a
// single error using simple values.
func ResponseError(statusCode int, title string, detail string) events.APIGatewayProxyResponse {
	jsonError := jsonapi.ErrorObject{
		Title:  title,
		Detail: detail,
		Status: fmt.Sprint(statusCode),
	}
	return ResponseErrors(statusCode, []*jsonapi.ErrorObject{&jsonError})
}

// ResponseSuccess is a factory function to generate a successful JSON API
// response with a given payload.
func ResponseSuccess(payload interface{}) events.APIGatewayProxyResponse {
	buf := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(buf, payload)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503, Body: err.Error()}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       buf.String(),
	}
}
