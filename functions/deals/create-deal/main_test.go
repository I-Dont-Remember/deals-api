package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

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

func Test_createDeal(t *testing.T) {
	// TODO: shouldn't need this setup, ideally we are only passing a subset of information in the body anyway
	d := models.Deal{
		ID:          "new-deal-id",
		Description: "description of deal",
	}

	jsonStr, err := json.Marshal(d)
	if err != nil {
		log.Print("Failed setting up test")
		os.Exit(3)
	}

	tests := []requestTest{
		{
			description: "200 and finds the correct item",
			request: events.APIGatewayProxyRequest{
				Body: string(jsonStr),
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
		response, err := createDeal(test.request, dbClient)
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
