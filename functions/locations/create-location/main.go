package main

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type locationBody struct {
	Name           string `json:"name"`
	DisplayAddress string `json:"address"`
}

func createLocation(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	if err := helpers.AuthMiddleware(request); err != nil {
		return helpers.ErrResponse("Failed authenticating", err, http.StatusUnauthorized)
	}

	slug := request.PathParameters["slug"]
	campus, err := db.GetCampus(slug)
	if err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	// if couldn't find campus matching path param
	if campus.Slug == "" {
		return helpers.ErrResponse("Bad request", err, http.StatusBadRequest)
	}

	body := locationBody{}
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	newID, err := uuid.NewV4()
	if err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	// TODO: find a better way to do input validaton than writing all our own checks
	if body.Name == "" {
		return helpers.ErrResponse("Need a location name", nil, http.StatusBadRequest)
	}

	location := models.Location{
		ID:         newID.String(),
		Name:       body.Name,
		CampusSlug: slug,
		Deals:      []string{},
	}

	if body.DisplayAddress != "" {
		location.DisplayAddress = body.DisplayAddress
	}

	err = db.CreateLocation(location)
	if err != nil {
		return helpers.ErrResponse("Issue creating location", err, http.StatusInternalServerError)
	}

	campus.Locations = append(campus.Locations, location.ID)
	_, err = db.UpdateCampus(campus)
	if err != nil {
		return helpers.ErrResponse("Failed updating campus", err, http.StatusInternalServerError)
	}

	marshalled, err := json.Marshal(location)
	if err != nil {
		return helpers.ErrResponse("Failed marshalling location", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusCreated)
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
