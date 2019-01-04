package helpers

import (
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// AuthMiddleware returns an error if it could not authorize the request
func AuthMiddleware(request events.APIGatewayProxyRequest) error {
	headers := request.Headers

	// log.Print(headers)
	val, ok := headers["X-Dot-Auth"]
	if ok {
		if val == os.Getenv("API_AUTH") {
			return nil
		}
	}
	return errors.New("Failed to authorize")
}
