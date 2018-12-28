package main

import (
	"errors"
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description string
	request     events.APIGatewayProxyRequest
	expect      string
	err         error
}

func Test_getCampuses(t *testing.T) {
	// Need to mock dynamodb values, since we already can pass the correct apigateway requests

	tests := []requestTest{
		{
			description: "",
			request: events.APIGatewayProxyRequest{
				Body: "",
				QueryStringParameters: nil,
			},
			expect: "",
			err:    nil,
		},
	}

	for _, test := range tests {
		mockClient := db.Mock{
			GetCampusesFunc: func() ([]models.Campus, error) {
				return []models.Campus{models.Campus{Slug: "hoopy-scoopy"}}, errors.New("ruh roh issue")
			},
		}

		response, err := getCampuses(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
