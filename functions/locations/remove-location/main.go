package main

import (
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func removeLocation(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	err := db.RemoveLocation(id)
	if err != nil {
		return helpers.ErrResponse("Issue removing location", err, http.StatusFailedDependency)
	}

	return helpers.Response("Location removed", http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something - look at net/http handlers setup ----> golden ticket
	dbClient, _ := db.Connect()
	return removeLocation(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
