package main

import (
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return helpers.Response("Not Implemented", http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
