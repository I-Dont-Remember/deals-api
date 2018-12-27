package main

import (
	"encoding/json"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createDeal(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	deal := models.Deal{}

	// TODO: forces the user to create the json exactly as the object, is that what we want? Shouldn't we handle creating id's and such?
	err := json.Unmarshal([]byte(request.Body), &deal)
	if err != nil {
		return helpers.ErrResponse("Invalid deal item", err, http.StatusBadRequest)
	}

	err = db.CreateDeal(deal)
	if err != nil {
		return helpers.ErrResponse("Issue creating deal", err, http.StatusFailedDependency)
	}

	marshalled, err := json.Marshal(deal)
	if err != nil {
		return helpers.ErrResponse("Failed marshalling deal", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something - look at net/http handlers setup ----> golden ticket
	dbClient, _ := db.Connect()
	return createDeal(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
