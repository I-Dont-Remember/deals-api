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

func getCampuses(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	campuses, err := db.GetCampuses()
	if err != nil {
		return helpers.ErrResponse("Issue getting campuses", err, http.StatusFailedDependency)
	}

	marshalled, err := json.Marshal(campuses)
	if err != nil {
		log.Print("Error marshalling campuses...")
		return helpers.ErrResponse("Failed marshalling campuses", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: handle errors and possibly extract this code with function pointers or something
	dbClient, _ := db.Connect()
	return getCampuses(request, dbClient)
}

func main() {
	lambda.Start(Handler)
}
