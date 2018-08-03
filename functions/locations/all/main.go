package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getLocations() {

}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	locations := getLocations()
	marshalled, err := json.Marshal(locations)
	if err != nil {
		log.Print("Error marshalling deals...")
		os.Exit(2)
	}

	// These are necessary alongside API Gateway CORS enabling
	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}

	return events.APIGatewayProxyResponse{
		Body:       string(marshalled),
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
