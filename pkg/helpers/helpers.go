package helpers

// most of this file adapted from https://ewanvalentine.io/serverless-start-ups-in-golang-part-1/
import (
	"encoding/json"
	"log"
	"os"

	"github.com/I-Dont-Remember/deals-api/pkg/db"

	"github.com/aws/aws-lambda-go/events"
)

// DbSetupForTest handles checking local vs test mode
func DbSetupForTest(mockClient db.Mock) db.DB {
	var dbClient db.DB
	var err error

	// assume 'test' env unless explicitly 'local'
	dbClient = mockClient
	if os.Getenv("API_ENV") == "local" {
		log.Print("Running tests in local mode")
		dbClient, err = db.Connect()
		if err != nil {
			log.Print("Unable to setup DB client in test")
			os.Exit(1)
		}
	}

	return dbClient
}

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
