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
	tableName = "Deals"
)

// Deal is a rough analogue to the DynamoDB schema; json for easy adding to table
type Deal struct {
	Id       string `json:"Id"`
	Day      string `json:"Day"`
	Location string `json:"Location"`
	Deal     string `json:"Deal"`
}

func buildParams() *dynamodb.ScanInput {
	filt := expression.Name("Id").AttributeExists()
	proj := expression.NamesList(expression.Name("Id"),
		expression.Name("Day"),
		expression.Name("Location"),
		expression.Name("Deal"))

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

func init() {
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
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Fetching Deals...")
	var deals []Deal

	for _, i := range output.Items {
		deal := Deal{}

		err := dynamodbattribute.UnmarshalMap(i, &deal)

		if err != nil {
			log.Print("Error unmarshaling item...")
			os.Exit(2)
		}

		log.Print(deal)
		deals = append(deals, deal)
	}

	marshalled, err := json.Marshal(deals)
	if err != nil {
		log.Print("Error marshalling deals...")
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
