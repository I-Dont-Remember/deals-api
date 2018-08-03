package main

import (
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description string
	request     events.APIGatewayProxyRequest
	expect      string
	err         error
}

// type APIGatewayProxyRequest struct {
//     Resource              string                        `json:"resource"` // The resource path defined in API Gateway
//     Path                  string                        `json:"path"`     // The url path for the caller
//     HTTPMethod            string                        `json:"httpMethod"`
//     Headers               map[string]string             `json:"headers"`
//     QueryStringParameters map[string]string             `json:"queryStringParameters"`
//     PathParameters        map[string]string             `json:"pathParameters"`
//     StageVariables        map[string]string             `json:"stageVariables"`
//     RequestContext        APIGatewayProxyRequestContext `json:"requestContext"`
//     Body                  string                        `json:"body"`
//     IsBase64Encoded       bool                          `json:"isBase64Encoded,omitempty"`
// }

func TestHandler(t *testing.T) {
	// Need to mock dynamodb values, since we already can pass the correct apigateway requests

	tests := []requestTest{
		{
			description: "Doesn't allow request without correct secret",
			request: events.APIGatewayProxyRequest{
				Body: `{"superSecret":"idontwork","key":"value"}`,
				QueryStringParameters: nil,
			},
			expect: "Forbidden",
			err:    nil,
		},
		{
			description: "Inputs item correctly",
			request: events.APIGatewayProxyRequest{
				Body: `{"superSecret":"HelloFreshOle","id":"1234","deal":"$4 food somewhere"}`,
				QueryStringParameters: nil,
			},
			expect: "",
			err:    nil,
		},
		// {
		// 	description: "",
		// 	request: events.APIGatewayProxyRequest{
		// 		Body: "",
		// 		QueryStringParameters: nil,
		// 	},
		// 	expect: "",
		// 	err:    nil,
		// },
	}

	for _, test := range tests {
		log.Print(test.description)
		response, err := Handler(test.request)
		assert.Nil(t, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
