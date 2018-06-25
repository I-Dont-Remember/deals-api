package main

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description string
	request     events.APIGatewayProxyRequest
	expect      string
	err         error
}

type dealsFilterTest struct {
	description string
	filter      string
}

// type deal struct {
// 	ID       string   `json:"id"`
// 	Location string   `json:"location"`
// 	Info     string   `json:"info"`
// 	Days     []string `json:"days"`
// }
func mockDynamo() []deal {
	return []deal{
		deal{"1", "abcd", "great deal", []string{"M", "Tu"}},
		deal{"2", "abcd", "great deal", []string{"M", "Tu", "W", "F"}},
		deal{"3", "abcd", "great deal", []string{"F"}},
		deal{"4", "abcd", "great deal", []string{"W", "F"}},
		deal{"5", "abcd", "great deal", []string{"M", "Tu"}},
		deal{"6", "abcd", "great deal", []string{"Sa", "Su"}},
		deal{"7", "cdef", "great deal", []string{"Th", "F", "Sa"}},
		deal{"8", "cdef", "great deal", []string{"F", "Sa"}},
		deal{"9", "cdef", "great deal", []string{"W", "Su"}},
		deal{"10", "cdef", "great deal", []string{"Tu", "Th"}},
	}
}
func Test_filterByLocation(t *testing.T) {
	tests := []dealsFilterTest{
		{"All items should match location filter", "abcd"},
	}

	for _, test := range tests {
		output := filterByLocation(mockDynamo(), test.filter)
		for _, item := range output {
			assert.Equal(t, test.filter, item.Location, test.description)
		}
	}
}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}

func Test_filterByDays(t *testing.T) {
	tests := []dealsFilterTest{
		{"Should only be Monday", "M"},
		{"Should be all days", "M,Tu,W,Th,F,Sa,Su"},
		{"Should have Tuesday & Thursday", "Tu,Th"},
	}

	for _, test := range tests {
		output := filterByDays(mockDynamo(), test.filter)
		// TODO: this is ugly, fix this to be less disgusting
		for _, item := range output {
			filterDays := strings.Split(test.filter, ",")
			passed := false
			for _, day := range filterDays {
				if contains(item.Days, day) {
					passed = true
				}
			}

			if !passed {
				t.Error("Incorrectly filtered days")
			}

		}
	}
}

// type APIGatewayProxyRequest struct {
//     Resource              string                        `json:"resource"` // The resource path defined in API Gateway
//     Path                  string                        `json:"path"`     // The url path for the caller
//     HTTPMethod            string                        `json:"httpMethod"`
//     Headers               map[string]string             `json:"headers"`
//     QueryStringParameters map[string]string             `json:"queryStringParameters"`
//     PathParameters        map[string]string             `json:"pathParameters"`
//     StageVariables        map[string]string             `json:"stageVariables"`
//     RequestContext        APIGatewayProxyRequestContext `json:"requestContext"`
//     Body                  string                        `json:"body"`
//     IsBase64Encoded       bool                          `json:"isBase64Encoded,omitempty"`
// }

func TestHandler(t *testing.T) {
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
		// {
		// 	description: "",
		// 	request: events.APIGatewayProxyRequest{
		// 		Body: "",
		// 		QueryStringParameters: nil,
		// 	},
		// 	expect: "",
		// 	err:    nil,
		// },
		// {
		// 	description: "",
		// 	request: events.APIGatewayProxyRequest{
		// 		Body: "",
		// 		QueryStringParameters: nil,
		// 	},
		// 	expect: "",
		// 	err:    nil,
		// },
	}

	os.Setenv("LAMBDA_ENV", "TEST")
	os.Setenv("TEST_DB", "TESTDB")
	for _, test := range tests {
		response, err := Handler(test.request)
		if err == nil {
			log.Print(response)
		}
		assert.NotEqual(t, test.expect, response.Body)
	}
}
