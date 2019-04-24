package email

import (
	"fmt"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
)

// Receive emails on AWS from IFTTT; keep db so that we don't break adjust function in local.go (this is bad practice)
func Receive(request events.APIGatewayProxyRequest, d db.DB) (events.APIGatewayProxyResponse, error) {
	fmt.Println(request.Body)
	return helpers.Response("", http.StatusOK)
}
