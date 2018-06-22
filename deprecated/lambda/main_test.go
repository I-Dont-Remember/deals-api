package main

import (
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	request events.APIGatewayProxyRequest
	expect  string
	err     error
}

func TestHandler(t *testing.T) {
	tests := []Test{
		{
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "",
			err:     nil,
		},
		{
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		if err == nil {
			log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
