package main

import (
	"encoding/json"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getLocation(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	location, err := db.GetLocation(id)
	if err != nil {
		return helpers.ErrResponse("Issue getting location", err, http.StatusFailedDependency)
	}

	// return 404 because if you wanted a specific item and it's not there, a 200 just makes for more callee error handling
	if location.ID == "" {
		return helpers.ErrResponse("location not found", nil, 404)
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
	return getLocation(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
