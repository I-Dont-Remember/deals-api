package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getDeals(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	deals, err := db.GetDeals()
	if err != nil {
		return helpers.ErrResponse("Issue getting deals", err, http.StatusFailedDependency)
	}

	marshalled, err := json.Marshal(deals)
	if err != nil {
		log.Print("Error marshalling deals...")
		return helpers.ErrResponse("Failed marshalling deals", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something - look at net/http handlers setup ----> golden ticket
	dbClient, _ := db.Connect()
	return getDeals(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
