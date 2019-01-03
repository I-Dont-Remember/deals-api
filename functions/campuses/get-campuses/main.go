package main

import (
	"github.com/I-Dont-Remember/deals-api/functions/campuses"
	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something
	dbClient, _ := db.Connect()
	return campuses.Get(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
