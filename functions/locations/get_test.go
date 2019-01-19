package locations

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {

	tests := []helpers.RequestTest{}

	rt := helpers.NewRequestTest()
	rt.Description = "200 got the locations"
	rt.ExpectedStatus = 200
	tests = append(tests, rt)

	rt = helpers.NewRequestTest()
	rt.Description = "200 got the expanded locations"
	rt.Request.QueryStringParameters = map[string]string{"expand": "deals"}
	rt.ExpectedStatus = 200
	tests = append(tests, rt)

	for _, test := range tests {

		test.Setup()

		response, _ := Get(test.Request, helpers.DbSetupForTest(test.MockClient))

		log.Print(response)

		assert.Equal(t, test.ExpectedStatus, response.StatusCode)
	}
}
