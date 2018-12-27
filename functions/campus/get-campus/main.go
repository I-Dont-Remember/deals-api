package main

import (
	"encoding/json"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getCampus(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	slug := request.PathParameters["slug"]
	campus, err := db.GetCampus(slug)
	if err != nil {
		return helpers.ErrResponse("Issue getting campus", err, http.StatusFailedDependency)
	}

	// return 404 because if you wanted a specific item and it's not there, a 200 just makes for more callee error handling
	if campus.Slug == "" {
		return helpers.ErrResponse("Campus not found", nil, 404)
	}

	marshalled, err := json.Marshal(campus)
	if err != nil {
		return helpers.ErrResponse("Failed marshalling campus", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something - look at net/http handlers setup ----> golden ticket
	dbClient, _ := db.Connect()
	return getCampus(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
