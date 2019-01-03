package campuses

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
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

func Test_GetOne(t *testing.T) {
	// Need to mock dynamodb values, since we already can pass the correct apigateway requests

	tests := []requestTest{
		{
			description: "200 and finds the correct item",
			request: events.APIGatewayProxyRequest{
				Body: "",
				PathParameters: map[string]string{
					"slug": "iowa-city",
				},
			},
			expect: "",
			err:    nil,
		},
		{
			description: "404 on not finding an item",
			request: events.APIGatewayProxyRequest{
				Body: "",
				PathParameters: map[string]string{
					"slug": "iowa-cityies",
				},
			},
			expect: "",
			err:    nil,
		},
	}

	for _, test := range tests {
		mockClient := db.Mock{
			GetCampusFunc: func(slug string) (models.Campus, error) {
				return models.Campus{}, nil
			},
		}
		response, err := GetOne(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
