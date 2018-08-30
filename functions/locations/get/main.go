package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	region    = "us-east-2"
	tableName = "devLocations"
)

type location struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Campus  string `json:"campus"`
	Name    string `json:"name"`
}

func buildParams() *dynamodb.ScanInput {
	filt := expression.Name("id").AttributeExists()
	proj := expression.NamesList(expression.Name("id"),
		expression.Name("address"),
		expression.Name("campus"),
		expression.Name("name"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Print("Error building expression...")
		os.Exit(2)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	return params
}

func getLocations() ([]location, error) {
	var locations []location
	svc, err := db.CreateConnection(region)
	if err != nil {
		return nil, err
	}

	params := buildParams()

	result, err := svc.Scan(params)
	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		l := location{}
		err := dynamodbattribute.UnmarshalMap(i, &l)
		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}

	return locations, nil
}

// Handler is lambda handler
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	locations, err := getLocations()
	if err != nil {
		return helpers.ErrResponse("", err, http.StatusInternalServerError)
	}
	marshalled, err := json.Marshal(locations)
	if err != nil {
		log.Print("Error marshalling deals...")
		return helpers.ErrResponse("", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
