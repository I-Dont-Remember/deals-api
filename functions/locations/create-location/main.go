package main

import (
	"encoding/json"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createLocation(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	location := models.Location{}

	// TODO: forces the user to create the json exactly as the object, is that what we want? Shouldn't we handle creating id's and such?
	err := json.Unmarshal([]byte(request.Body), &location)
	if err != nil {
		return helpers.ErrResponse("Invalid Location item", err, http.StatusBadRequest)
	}

	err = db.CreateLocation(location)
	if err != nil {
		return helpers.ErrResponse("Issue creating location", err, http.StatusFailedDependency)
	}

	marshalled, err := json.Marshal(location)
	if err != nil {
		return helpers.ErrResponse("Failed marshalling location", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something - look at net/http handlers setup ----> golden ticket
	dbClient, _ := db.Connect()
	return createLocation(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
