package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	region    = "us-east-2"
	output    *dynamodb.ScanOutput
	tableName = "Locations"
)

// Location is a schema (ish) for the Location table
type Location struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

func buildParams() *dynamodb.ScanInput {
	filt := expression.Name("ID").AttributeExists()
	proj := expression.NamesList(expression.Name("ID"),
		expression.Name("Name"))

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

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Fetching Locations...")
	var locations []Location

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Print("Error getting session")
		log.Print(err)
		// this might not be good inside lambda function???
		os.Exit(2)
	}

	svc := dynamodb.New(sess)
	params := buildParams()

	result, err := svc.Scan(params)
	if err != nil {
		log.Print("Error scanning...")
		os.Exit(2)
	}

	output = result

	for _, i := range output.Items {
		location := Location{}

		err := dynamodbattribute.UnmarshalMap(i, &location)

		if err != nil {
			log.Print("Error unmarshaling item...")
			os.Exit(2)
		}

		log.Print(location)
		locations = append(locations, location)
	}

	marshalled, err := json.Marshal(locations)
	if err != nil {
		log.Print("Error marshalling locations...")
		os.Exit(2)
	}

	// These are necessary alongside API Gateway CORS enabling
	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}

	return events.APIGatewayProxyResponse{
		Body:       string(marshalled),
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
