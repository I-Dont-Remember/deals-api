package campuses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
)

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
