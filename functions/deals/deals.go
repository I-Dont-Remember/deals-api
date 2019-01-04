package deals

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
)

type dealBody struct {
	Days        []string `json:"days"`
	Description string   `json:"description"`
	AllDay      bool     `json:"all_day"`
	Types       []string `json:"types"`
}

// Get fetches all deals
func Get(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
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

// Remove removes a deal
func Remove(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	err := db.RemoveDeal(id)
	if err != nil {
		return helpers.ErrResponse("Issue removing deal", err, http.StatusFailedDependency)
	}

	return helpers.Response("Deal removed", http.StatusOK)
}

// Create makes a new deal
// TODO: this should be an idempotent function, we need a way to make sure the tables aren't ever out of sync
func Create(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	locationID := request.PathParameters["location-id"]

	// fetch the location and check it exists
	location, err := db.GetLocation(locationID)
	if err != nil {
		return helpers.ErrResponse("Internal Error", err, http.StatusInternalServerError)
	}

	if location.ID == "" {
		return helpers.Response("Bad ID", http.StatusBadRequest)
	}

	newID, err := uuid.NewV4()
	if err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	// TODO: need a better way to validate we've gotten all necessary info & it's gucci
	body := dealBody{}
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	// TODO: ignore enum validation for Days/Type for now :D

	// create deal
	deal := models.Deal{
		ID:          newID.String(),
		LocationID:  locationID,
		Description: body.Description,
		AllDay:      body.AllDay,
		Days:        body.Days,
		Types:       body.Types,
	}

	err = db.CreateDeal(deal)
	if err != nil {
		return helpers.ErrResponse("Issue creating deal", err, http.StatusFailedDependency)
	}

	location.Deals = append(location.Deals, deal.ID)
	_, err = db.UpdateLocation(location)
	if err != nil {
		return helpers.ErrResponse("Failed updating location", err, http.StatusInternalServerError)
	}

	marshalled, err := json.Marshal(deal)
	if err != nil {
		return helpers.ErrResponse("Failed marshalling deal", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusCreated)
}
