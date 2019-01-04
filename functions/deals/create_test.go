package deals

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description      string
	bodyMap          map[string]interface{}
	request          events.APIGatewayProxyRequest
	expectedStatus   int
	dbMockFunc       func(models.Deal) error
	locationMockFunc func(string) (models.Location, error)
	updateMockFunc   func(models.Location) (models.Location, error)
}

func Test_Createl(t *testing.T) {
	tests := []requestTest{
		{
			description: "201 creates a Deal",
			bodyMap: map[string]interface{}{
				"description": "a deal description",
				"all_day":     true,
				"types":       []string{"Event"},
				"days":        []string{"Mon", "Wed", "Fri"},
			},
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"location-id": "abcdefghighs"},
			},
			expectedStatus: 201,
			dbMockFunc: func(models.Deal) error {
				return nil
			},
			locationMockFunc: func(id string) (models.Location, error) {
				return models.Location{ID: "abcddefegegeeg"}, nil
			},
			updateMockFunc: func(models.Location) (models.Location, error) {
				return models.Location{}, nil
			},
		},
		{
			description: "400 if given invalid id",
			bodyMap:     map[string]interface{}{},
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"location-id": ""},
			},
			expectedStatus: 400,
			dbMockFunc: func(models.Deal) error {
				return nil
			},
			locationMockFunc: func(id string) (models.Location, error) {
				return models.Location{ID: ""}, nil
			},
			updateMockFunc: func(models.Location) (models.Location, error) {
				return models.Location{}, nil
			},
		},
	}

	for _, test := range tests {
		fmt.Println("[*] " + test.description)
		bytes, err := json.Marshal(test.bodyMap)
		if err != nil {
			log.Print("Failed setting up for test")
			os.Exit(3)
		}
		test.request.Body = string(bytes)

		mockClient := db.Mock{
			CreateDealFunc:     test.dbMockFunc,
			GetLocationFunc:    test.locationMockFunc,
			UpdateLocationFunc: test.updateMockFunc,
		}
		response, err := Create(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)

		assert.Equal(t, test.expectedStatus, response.StatusCode)
	}
}
