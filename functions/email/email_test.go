package email

import (
	"log"
	"testing"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/stretchr/testify/assert"
)

func Test_Receive(t *testing.T) {

	tests := []helpers.RequestTest{}

	rt := helpers.NewRequestTest()
	rt.Description = "201 posted the search data"
	rt.BodyMap = map[string]string{"subject": "hello world"}
	rt.ExpectedStatus = 200

	tests = append(tests, rt)

	for _, test := range tests {

		test.Setup()

		response, _ := Receive(test.Request, db.Dynamo{})

		log.Print(response)

		assert.Equal(t, test.ExpectedStatus, response.StatusCode)
	}
}
