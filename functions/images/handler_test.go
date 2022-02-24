package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sapslaj/aquapi/internal/dev"
)

func stringSliceContains(input []string, match string) bool {
	for _, s := range input {
		if s == match {
			return true
		}
	}
	return false
}

func useMockRequest(t *testing.T, r *events.APIGatewayProxyRequest) (*dev.MockRequest, map[string]interface{}, string) {
	mr := dev.NewMockRequest(handler, r)
	body, err := mr.GetBody()
	if err != nil {
		t.Fatal("error returned getting body", err)
	}
	bodystring, err := mr.GetBodyString()
	if err != nil {
		t.Fatal("error returned getting body as a string", err)
	}
	return mr, body, bodystring
}

func TestReturnsOne(t *testing.T) {
	_, body, bodystring := useMockRequest(t, &events.APIGatewayProxyRequest{
		Resource:   "/",
		Path:       "/",
		HTTPMethod: "GET",
	})
	data := body["data"].([]interface{})
	if len(data) != 1 {
		t.Fatal("data length is not 1: ", body, bodystring)
	}
}

func TestReturnsCorrectCount(t *testing.T) {
	_, body, bodystring := useMockRequest(t, &events.APIGatewayProxyRequest{
		Resource:   "/",
		Path:       "/",
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"count": "5",
		},
	})
	data := body["data"].([]interface{})
	if len(data) != 5 {
		t.Fatal("data length is not 1: ", body, bodystring)
	}
}

func TestErrorsOnInvalidCount(t *testing.T) {
	_, body, bodystring := useMockRequest(t, &events.APIGatewayProxyRequest{
		Resource:   "/",
		Path:       "/",
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"count": "invalid",
		},
	})
	errors := body["errors"].([]interface{})
	if len(errors) == 0 {
		t.Fatal("errors length is 0: ", body, bodystring)
	}
}

func TestErrorsOnLargeCount(t *testing.T) {
	_, body, bodystring := useMockRequest(t, &events.APIGatewayProxyRequest{
		Resource:   "/",
		Path:       "/",
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"count": "99999999",
		},
	})
	errors := body["errors"].([]interface{})
	if len(errors) == 0 {
		t.Fatal("errors length is 0: ", body, bodystring)
	}
}

func TestAcceptsNsfwParam(t *testing.T) {
	t.Skip("test is a lil bonkers right now")
	_, body, bodystring := useMockRequest(t, &events.APIGatewayProxyRequest{
		Resource:   "/",
		Path:       "/",
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"nsfw": "only",
		},
	})
	data := body["data"].([]interface{})
	for i, d := range data {
		datum := d.(map[string]interface{})
		attributes := datum["attributes"].(map[string][]string)
		if !stringSliceContains(attributes["tags"], "nsfw") {
			t.Fatal("tags for "+fmt.Sprint(i)+" does not contain nsfw", body, bodystring)
		}
	}
}
