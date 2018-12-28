package main

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description string
	request     events.APIGatewayProxyRequest
	expect      string
	err         error
}

func Test_removeDeal(t *testing.T) {
	tests := []requestTest{
		{
			description: "200 and finds the correct item",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"id": "new-deal-id",
				},
			},
			expect: "",
			err:    nil,
		},
	}

	for _, test := range tests {
		mockClient := db.Mock{
			RemoveDealFunc: func(id string) error {
				return nil
			},
		}
		response, err := removeDeal(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
