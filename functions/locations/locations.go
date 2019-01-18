package locations

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
)

type locationBody struct {
	Name           string `json:"name"`
	DisplayAddress string `json:"address"`
}

// Create makes a new location
func Create(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	if err := helpers.AuthMiddleware(request); err != nil {
		return helpers.ErrResponse("Failed authenticating", err, http.StatusUnauthorized)
	}

	slug := request.PathParameters["slug"]
	campus, err := db.GetCampus(slug)
	if err != nil {
		return helpers.ErrResponse("Error getting campus", err, http.StatusInternalServerError)
	}

	// TODO: should have some way to check that we aren't doubling up a location

	// if couldn't find campus matching path param
	if campus.Slug == "" {
		return helpers.ErrResponse("Bad request", err, http.StatusBadRequest)
	}

	body := locationBody{}
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return helpers.ErrResponse("Error getting body", err, http.StatusInternalServerError)
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

// GetOne fetches a single location
func GetOne(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["location-id"]
	log.Println("location GetOne id: " + id)
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

// Get retrieves all locations
func Get(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	locations, err := db.GetLocations()
	if err != nil {
		return helpers.ErrResponse("Issue getting locations", err, http.StatusFailedDependency)
	}

	marshalled, err := json.Marshal(locations)
	if err != nil {
		log.Print("Error marshalling locations...")
		return helpers.ErrResponse("Failed marshalling locations", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

// Remove gets rid of location
func Remove(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["location-id"]
	err := db.RemoveLocation(id)
	if err != nil {
		return helpers.ErrResponse("Issue removing location", err, http.StatusFailedDependency)
	}

	return helpers.Response("Location removed", http.StatusOK)
}
