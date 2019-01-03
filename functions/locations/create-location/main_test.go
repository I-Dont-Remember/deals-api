package main

import (
	"encoding/json"
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
	description    string
	bodyMap        map[string]string
	request        events.APIGatewayProxyRequest
	expectedStatus int
	authValue      string
	dbMockFunc     func(models.Location) error
	campusMockFunc func(string) (models.Campus, error)
	UpdateMockFunc func(models.Campus) (models.Campus, error)
}

func Test_createLocation(t *testing.T) {
	tests := []requestTest{
		{
			description: "201 created the Location",
			bodyMap:     map[string]string{"name": "new-location"},
			request: events.APIGatewayProxyRequest{
				Headers:        map[string]string{"x-dot-auth": "success"},
				PathParameters: map[string]string{"slug": "madison-wi"},
			},
			authValue:      "success",
			expectedStatus: 201,
			dbMockFunc: func(models.Location) error {
				return nil
			},
			campusMockFunc: func(slug string) (models.Campus, error) {
				return models.Campus{Slug: "madison-wi"}, nil
			},
			UpdateMockFunc: func(campus models.Campus) (models.Campus, error) {
				return models.Campus{}, nil
			},
		},
		{
			description: "400 if not given a name",
			bodyMap:     map[string]string{"whatsit": "not-a-name"},
			request: events.APIGatewayProxyRequest{
				Headers:        map[string]string{"x-dot-auth": "success"},
				PathParameters: map[string]string{"slug": "madison-wi"},
			},
			authValue:      "success",
			expectedStatus: 400,
			dbMockFunc: func(models.Location) error {
				return nil
			},
			campusMockFunc: func(slug string) (models.Campus, error) {
				return models.Campus{Slug: "madison-wi"}, nil
			},
			UpdateMockFunc: func(campus models.Campus) (models.Campus, error) {
				return models.Campus{}, nil
			},
		},
		{
			description: "400 if given an unknown slug",
			bodyMap:     map[string]string{"whatsit": "not-a-name"},
			request: events.APIGatewayProxyRequest{
				Headers:        map[string]string{"x-dot-auth": "success"},
				PathParameters: map[string]string{"slug": "unknown-location"},
			},
			authValue:      "success",
			expectedStatus: 400,
			dbMockFunc: func(models.Location) error {
				return nil
			},
			campusMockFunc: func(slug string) (models.Campus, error) {
				return models.Campus{}, nil
			},
			UpdateMockFunc: func(campus models.Campus) (models.Campus, error) {
				return models.Campus{}, nil
			},
		},
		{
			description: "401 Prevents unauthorized access",
			bodyMap:     map[string]string{"whatsit": "not-a-name"},
			request: events.APIGatewayProxyRequest{
				Headers:        map[string]string{"x-dot-auth": "dont-let-me-in"},
				PathParameters: map[string]string{"slug": "unknown-location"},
			},
			authValue:      "failure-test",
			expectedStatus: 401,
			dbMockFunc: func(models.Location) error {
				return nil
			},
			campusMockFunc: func(slug string) (models.Campus, error) {
				return models.Campus{}, nil
			},
			UpdateMockFunc: func(campus models.Campus) (models.Campus, error) {
				return models.Campus{}, nil
			},
		},
	}

	for _, test := range tests {
		os.Setenv("API_AUTH", test.authValue)

		bytes, err := json.Marshal(test.bodyMap)
		if err != nil {
			log.Print("Failed setting up for test")
			os.Exit(3)
		}
		test.request.Body = string(bytes)

		mockClient := db.Mock{
			CreateLocationFunc: test.dbMockFunc,
			GetCampusFunc:      test.campusMockFunc,
			UpdateCampusFunc:   test.UpdateMockFunc,
		}
		response, err := createLocation(test.request, helpers.DbSetupForTest(mockClient))
		log.Print(response)

		assert.Equal(t, test.expectedStatus, response.StatusCode)
	}
}
