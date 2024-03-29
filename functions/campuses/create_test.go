package campuses

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
	rt.Description = "201 created the Campus"
	rt.BodyMap = map[string]string{"slug": "uw-madison", "display_name": "University of Wisconsin-Madison"}
	rt.MockClient.GetCampusFunc = func(slug string) (models.Campus, error) {
		return models.Campus{}, nil
	}
	rt.ExpectedStatus = 201
	tests = append(tests, rt)

	rt = helpers.NewRequestTest()
	rt.Description = "400 if not given a slug"
	rt.BodyMap = map[string]string{"whatsit": "not-a-slug"}
	rt.ExpectedStatus = 400
	tests = append(tests, rt)

	// slug already exists
	rt = helpers.NewRequestTest()
	rt.Description = "409 if slug already exists"
	rt.BodyMap = map[string]string{"slug": "i-already-exist"}
	rt.MockClient.GetCampusFunc = func(slug string) (models.Campus, error) {
		return models.Campus{Slug: "i-already-exist"}, nil
	}
	rt.ExpectedStatus = 409
	tests = append(tests, rt)

	rt = helpers.NewRequestTest()
	rt.Description = "401 Prevents unauthorized access"
	rt.BodyMap = map[string]string{"whatsit": "not-a-slug"}
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
