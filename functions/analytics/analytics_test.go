package analytics

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/stretchr/testify/assert"
)

func Test_Analytics(t *testing.T) {

	tests := []helpers.RequestTest{}

	rt := helpers.NewRequestTest()
	rt.Description = "201 posted the search data"
	rt.BodyMap = map[string]string{"timestamp": "1548118356297", "searchTerm": "wings"}
	rt.Request.PathParameters = map[string]string{"slug": "madison-wi"}
	rt.ExpectedStatus = 201
	rt.MockClient.InputSearchAnalyticsFunc = func(s models.SearchData) error {
		return nil
	}
	tests = append(tests, rt)

	for _, test := range tests {

		test.Setup()

		response, _ := Analytics(test.Request, helpers.DbSetupForTest(test.MockClient))

		log.Print(response)

		assert.Equal(t, test.ExpectedStatus, response.StatusCode)
	}
}
