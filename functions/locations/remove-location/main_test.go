package main

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description string
	request     events.APIGatewayProxyRequest
	expect      string
	err         error
}

func Test_removeLocation(t *testing.T) {
	tests := []requestTest{
		{
			description: "200 and finds the correct item",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"id": "new-location-id",
				},
			},
			expect: "",
			err:    nil,
		},
	}

	for _, test := range tests {
		// TODO: can have data structures in mock that get filled
		// by a test 'setup' function, then the interface functions just
		// access those
		dbClient, _ := db.Connect()
		response, err := removeLocation(test.request, dbClient)
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
