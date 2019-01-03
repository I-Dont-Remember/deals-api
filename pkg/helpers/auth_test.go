package helpers

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type requestTest struct {
	description string
	request     events.APIGatewayProxyRequest
	expectError bool
	authValue   string
}

func Test_createLocation(t *testing.T) {
	tests := []requestTest{
		{
			description: "Succeeds with matching header & env variable values",
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{"x-dot-auth": "success-test"},
			},
			authValue:   "success-test",
			expectError: false,
		},
		{
			description: "Prevents unauthorized access",
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{"x-dot-auth": "let-me-in"},
			},
			authValue:   "failure-test",
			expectError: true,
		},
	}

	for _, test := range tests {
		os.Setenv("API_AUTH", test.authValue)
		err := AuthMiddleware(test.request)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
