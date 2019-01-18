package campuses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"
	"github.com/aws/aws-lambda-go/events"
)

type campusBody struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"display_name"`
}

// Get fetches all campuses from the DB wrapper & returns a http response
func Get(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
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

// GetOne returns a single campus
func GetOne(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
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

// Create creates a campus
func Create(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {
	if err := helpers.AuthMiddleware(request); err != nil {
		return helpers.ErrResponse("Failed authenticating", err, http.StatusUnauthorized)
	}

	body := campusBody{}
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	if body.Slug == "" {
		return helpers.Response("no slug", http.StatusBadRequest)
	}

	// check if slug already exists
	campus, err := db.GetCampus(body.Slug)
	if err != nil {
		return helpers.ErrResponse("Internal error", err, http.StatusInternalServerError)
	}

	if campus.Slug != "" {
		return helpers.Response("slug exists", http.StatusConflict)
	}

	// TODO: make sure error is related to the slug being missing not other errors, this will probably be roped into better error handling all around

	campus = models.Campus{
		Slug:        body.Slug,
		DisplayName: body.DisplayName,
		Locations:   []string{},
	}

	err = db.CreateCampus(campus)
	if err != nil {
		return helpers.ErrResponse("Issue creating campus", err, http.StatusInternalServerError)
	}

	marshalled, err := json.Marshal(campus)
	if err != nil {
		return helpers.ErrResponse("Failed marshalling campus", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusCreated)
}
