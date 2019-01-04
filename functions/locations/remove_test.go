package locations

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type removeTest struct {
	description    string
	request        events.APIGatewayProxyRequest
	authValue      string
	expectedStatus int
	err            error
}

func Test_Remove(t *testing.T) {
	tests := []removeTest{
		{
			description: "200 and removes the correct item",
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{"x-dot-auth": "success"},
				PathParameters: map[string]string{
					"id": "new-location-id",
				},
			},
			authValue:      "success",
			expectedStatus: 200,
			err:            nil,
		},
		{
			description: "401 and finds the correct item",
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{"x-dot-auth": "dont-let-me-in"},
				PathParameters: map[string]string{
					"id": "new-location-id",
				},
			},
			expectedStatus: 401,
			err:            nil,
		},
	}

	for _, test := range tests {
		mockClient := db.Mock{
			RemoveLocationFunc: func(id string) error {
				return nil
			},
		}
		response, err := Remove(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.Equal(t, test.expectedStatus, response.StatusCode)
	}
}
