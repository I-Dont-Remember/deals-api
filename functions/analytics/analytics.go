package analytics

import (
	"encoding/json"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/aws/aws-lambda-go/events"
)

// Search handles basic analytics for users searching
func Analytics(request events.APIGatewayProxyRequest, db db.DB) (events.APIGatewayProxyResponse, error) {

	body := models.SearchData{}
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return helpers.ErrResponse("Error getting search body", err, http.StatusInternalServerError)
	}

	err := db.InputSearchAnalytics(body)
	if err != nil {
		return helpers.ErrResponse("Error inputting search", err, http.StatusInternalServerError)

	}

	return helpers.Response("posted search", http.StatusCreated)
}
