package deals

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// type requestTest struct {
// 	description string
// 	request     events.APIGatewayProxyRequest
// 	expect      string
// 	err         error
// }

func Test_Get(t *testing.T) {
	// Need to mock dynamodb values, since we already can pass the correct apigateway requests

	tests := []requestTest{
		{
			description: "",
			request: events.APIGatewayProxyRequest{
				Body: "",
				QueryStringParameters: nil,
			},
			expectedStatus: 200,
		},
	}

	for _, test := range tests {
		mockClient := db.Mock{
			GetDealsFunc: func() ([]models.Deal, error) {
				return []models.Deal{models.Deal{}}, nil
			},
		}
		response, err := Get(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)
		if err == nil {
			//log.Print(response)
		}
		assert.NotEqual(t, test.expectedStatus, response.StatusCode)
	}
}
