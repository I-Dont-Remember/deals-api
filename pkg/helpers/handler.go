package helpers

// most of this file adapted from https://ewanvalentine.io/serverless-start-ups-in-golang-part-1/
import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Response is a wrapper for our http response
func Response(data string, statusCode int) (events.APIGatewayProxyResponse, error) {
	// These are necessary alongside API Gateway CORS enabling
	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
	return events.APIGatewayProxyResponse{
		Body:       data,
		StatusCode: statusCode,
		Headers:    headers,
	}, nil
}

// ErrResponse returns an error in a specified format
func ErrResponse(msg string, err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	data := map[string]string{
		"msg": msg,
	}

	if err != nil {
		data["error"] = err.Error()
	}

	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}, err
}
