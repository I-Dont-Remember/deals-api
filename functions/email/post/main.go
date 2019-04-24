package main

import (
	email "github.com/I-Dont-Remember/deals-api/functions/email"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return email.Receive(request, nil)
}

func main() {
	lambda.Start(Handler)
}
