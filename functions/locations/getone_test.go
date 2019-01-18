package locations

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/stretchr/testify/assert"
)

func Test_GetOne(t *testing.T) {

	tests := []helpers.RequestTest{}

	rt := helpers.NewRequestTest()
	rt.Description = "200 got the Location"
	rt.Request.PathParameters = map[string]string{"id": "existing-id"}
	rt.ExpectedStatus = 200
	rt.MockClient.GetLocationFunc = func(id string) (models.Location, error) {
		return models.Location{ID: "existing-id"}, nil
	}
	tests = append(tests, rt)

	// default DB response should get this, because we check if returned location id is ""
	rt = helpers.NewRequestTest()
	rt.Description = "404 if item wasn't found"
	rt.Request.PathParameters = map[string]string{"id": "abcddeefgdgdgs"}
	rt.ExpectedStatus = 404
	tests = append(tests, rt)

	for _, test := range tests {

		test.Setup()

		response, _ := GetOne(test.Request, helpers.DbSetupForTest(test.MockClient))

		log.Print(response)

		assert.Equal(t, test.ExpectedStatus, response.StatusCode)
	}
}
