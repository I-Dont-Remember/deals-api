package helpers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/models"
	"github.com/aws/aws-lambda-go/events"
)

// RequestTest stores a variety of useful information for testing Lambda Request handlers.
// Should have all off the same functions available as the db.Mock
type RequestTest struct {
	Description    string
	BodyMap        map[string]string
	Request        events.APIGatewayProxyRequest
	ExpectedStatus int
	AuthValue      string
	MockClient     db.Mock
}

// NewRequestTest sets up defaults for all functions so we only have to edit the items that a specific test cares about at that.
// This includes setting the Auth header to be valid so it only gets tested when intended.  This approach lets us easily
// use defaults while allowing swapping of any of the items if a test needs to.
func NewRequestTest() RequestTest {
	return RequestTest{
		Description: "default description",
		BodyMap:     map[string]string{},
		Request: events.APIGatewayProxyRequest{
			Headers: map[string]string{"X-Dot-Auth": "success"},
		},
		ExpectedStatus: 0,
		AuthValue:      "success",
		MockClient: db.Mock{
			CreateCampusFunc: func(c models.Campus) error {
				return nil
			},
			RemoveCampusFunc: func(slug string) error {
				return nil
			},
			GetCampusesFunc: func() ([]models.Campus, error) {
				return []models.Campus{}, nil
			},
			GetCampusFunc: func(slug string) (models.Campus, error) {
				return models.Campus{}, nil
			},
			UpdateCampusFunc: func(c models.Campus) (models.Campus, error) {
				return models.Campus{}, nil
			},
			CreateLocationFunc: func(l models.Location) error {
				return nil
			},
			RemoveLocationFunc: func(id string) error {
				return nil
			},
			GetLocationsFunc: func() ([]models.Location, error) {
				return []models.Location{}, nil
			},
			GetLocationFunc: func(id string) (models.Location, error) {
				return models.Location{}, nil
			},
			UpdateLocationFunc: func(l models.Location) (models.Location, error) {
				return models.Location{}, nil
			},
			GetDealsFunc: func() ([]models.Deal, error) {
				return []models.Deal{}, nil
			},
			RemoveDealFunc: func(id string) error {
				return nil
			},
			CreateDealFunc: func(d models.Deal) error {
				return nil
			},
		},
	}
}

// Setup handles setting up any extra steps for a test with RequestTest
func (r *RequestTest) Setup() {
	os.Setenv("API_AUTH", r.AuthValue)

	bytes, err := json.Marshal(r.BodyMap)
	if err != nil {
		log.Print("[!] failed setting up for test - " + r.Description)
		os.Exit(3)
	}
	r.Request.Body = string(bytes)

	log.Print(r.Description)
}
