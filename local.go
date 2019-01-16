package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/I-Dont-Remember/deals-api/functions/deals"

	"github.com/I-Dont-Remember/deals-api/functions/locations"

	"github.com/I-Dont-Remember/deals-api/functions/campuses"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/I-Dont-Remember/deals-api/pkg/db"

	"github.com/aws/aws-lambda-go/events"
)

// Massage the echo framework request/response to match our AWS Lambda handlers
func adjust(fn func(events.APIGatewayProxyRequest, db.DB) (events.APIGatewayProxyResponse, error)) func(echo.Context) error {
	return func(c echo.Context) error {
		dbClient, _ := db.Connect()

		// TODO: validate that the string joining nonsense we're doing is actually working correctly
		headers := map[string]string{}
		for k, v := range c.Request().Header {
			headers[k] = strings.Join(v[:], ",")
		}

		paramNames := c.ParamNames()
		paramMap := map[string]string{}
		for _, name := range paramNames {
			paramMap[name] = c.Param(name)
		}

		queryMap := map[string][]string{}
		queryParams := map[string]string{}
		queryMap = c.QueryParams()
		for k, v := range queryMap {
			queryParams[k] = strings.Join(v[:], ",")
		}

		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			panic(err)
		}

		request := events.APIGatewayProxyRequest{
			HTTPMethod:            c.Request().Method,
			Headers:               headers,
			PathParameters:        paramMap,
			QueryStringParameters: queryParams,
			Body: string(body),
		}
		proxyResponse, _ := fn(request, dbClient)

		// TODO: check this is actually doing what we thing it is
		for k, v := range proxyResponse.Headers {
			c.Response().Header().Set(k, v)
		}
		return c.JSONBlob(proxyResponse.StatusCode, []byte(proxyResponse.Body))
	}
}

func main() {
	os.Setenv("API_ENV", "local")
	// Make sure to set the auth env variable as if we were deployed
	os.Setenv("API_AUTH", "local")
	port := ":4500"

	e := echo.New()

	// To see specific header, use ${header:foo} which will show foo's value
	// same for seeing cookie, query, and form
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} status:${status} latency:${latency_human} out:${bytes_out} bytes X-Dot-Auth:${header:X-Dot-Auth}\n",
	}))

	// !! These should all match exactly with the serverless.yml
	e.GET("/campuses", adjust(campuses.Get))
	e.GET("/campuses/:slug", adjust(campuses.GetOne))
	e.POST("/campuses", adjust(campuses.Create))

	e.GET("/campuses/:slug/locations", adjust(locations.Get))
	e.GET("/locations/:location-id", adjust(locations.GetOne))
	e.POST("/campuses/:slug/locations", adjust(locations.Create))
	e.DELETE("/locations/:location-id", adjust(locations.Remove))

	e.GET("/locations/:location-id/deals", adjust(deals.Get))
	e.POST("/locations/:location-id/deals", adjust(deals.Create))
	e.DELETE("/deals/:deal-id", adjust(deals.Remove))

	e.Logger.Fatal(e.Start(port))
}
