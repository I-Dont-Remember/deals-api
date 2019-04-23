package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// AuthMiddleware returns an error if it could not authorize the request
func AuthMiddleware(request events.APIGatewayProxyRequest) error {
	headers := request.Headers

	fmt.Printf("Checking headers %+v\n", headers)
	if val, ok := headers["X-Dot-Auth"]; ok {
		if val == os.Getenv("API_AUTH") {
			return nil
		}
	}

	if val, ok := headers["x-dot-auth"]; ok {
		if val == os.Getenv("API_AUTH") {
			return nil
		}
	}

	return errors.New("Failed to authorize")
}
