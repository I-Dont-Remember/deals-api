package main

import (
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetCampuses(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	return helpers.Response(db.GetCampuses()[0].Slug, http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dbClient, _ := db.Connect()
	return GetCampuses(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
