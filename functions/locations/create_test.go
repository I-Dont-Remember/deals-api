package locations

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {

	tests := []helpers.RequestTest{}

	rt := helpers.NewRequestTest()
	rt.Description = "201 created the Location"
	rt.BodyMap = map[string]string{"name": "new-location"}
	rt.Request.PathParameters = map[string]string{"slug": "madison-wi"}
	rt.ExpectedStatus = 201
	rt.MockClient.GetCampusFunc = func(slug string) (models.Campus, error) {
		return models.Campus{Slug: "madison-wi"}, nil
	}
	tests = append(tests, rt)

	rt = helpers.NewRequestTest()
	rt.Description = "400 if not given a name"
	rt.BodyMap = map[string]string{"whatsit": "not-a-name"}
	rt.Request.PathParameters = map[string]string{"slug": "madison-wi"}
	rt.ExpectedStatus = 400
	rt.MockClient.GetCampusFunc = func(slug string) (models.Campus, error) {
		return models.Campus{Slug: "madison-wi"}, nil
	}
	tests = append(tests, rt)

	rt = helpers.NewRequestTest()
	rt.Description = "400 if given an unknown slug"
	rt.BodyMap = map[string]string{"whatsit": "not-a-name"}
	rt.Request.PathParameters = map[string]string{"slug": "unknown-location"}
	rt.ExpectedStatus = 400
	tests = append(tests, rt)

	rt = helpers.NewRequestTest()
	rt.Description = "401 Prevents unauthorized access"
	rt.BodyMap = map[string]string{"whatsit": "not-a-name"}
	rt.Request.Headers = map[string]string{"X-Dot-Auth": "dont-let-me-in"}
	rt.AuthValue = "failure-test"
	rt.ExpectedStatus = 401
	tests = append(tests, rt)

	for _, test := range tests {

		test.Setup()

		response, _ := Create(test.Request, helpers.DbSetupForTest(test.MockClient))

		log.Print(response)

		assert.Equal(t, test.ExpectedStatus, response.StatusCode)
	}
}
