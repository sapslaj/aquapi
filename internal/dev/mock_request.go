package dev

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type MockRequest struct {
	Ctx      context.Context
	Handler  interface{}
	Request  *events.APIGatewayProxyRequest
	Response *events.APIGatewayProxyResponse
	Executed bool
}

func NewMockRequest(handler interface{}, request *events.APIGatewayProxyRequest) *MockRequest {
	return &MockRequest{
		Ctx:      context.TODO(),
		Handler:  handler,
		Request:  request,
		Executed: false,
	}
}

func (mr *MockRequest) Execute() error {
	response, err := mr.Handler.(func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error))(mr.Ctx, *mr.Request)
	mr.Response = &response
	if err != nil {
		return err
	}
	return nil
}

func (mr *MockRequest) GetResponse() (events.APIGatewayProxyResponse, error) {
	var err error
	if !mr.Executed {
		err = mr.Execute()
	}
	return *mr.Response, err
}

func (mr *MockRequest) GetBody() (map[string]interface{}, error) {
	response, err := mr.GetResponse()
	var body map[string]interface{}
	if err != nil {
		return body, err
	}
	err = json.Unmarshal([]byte(response.Body), &body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (mr *MockRequest) GetBodyString() (string, error) {
	body, err := mr.GetBody()
	if err != nil {
		return "", err
	}
	formattedBody, err := json.MarshalIndent(body, "", "  ")
	return string(formattedBody), err
}

func (mr *MockRequest) PrintResponse() {
	response, err := mr.GetResponse()
	fmt.Printf("response: %#v\n", response)
	if err != nil {
		fmt.Printf("err: %v\n\n", err)
	}

	body, err := mr.GetBodyString()
	fmt.Printf("body: %v\n", body)
	if err != nil {
		fmt.Printf("err: %v\n\n", err)
	}
}
