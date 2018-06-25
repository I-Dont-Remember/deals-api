package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	region        = "us-east-2"
	output        *dynamodb.ScanOutput
	tableName     = "devDeals"
	localEndpoint = "http://localhost:4569/"
)

// Deal is a rough analogue to the DynamoDB schema;  should move to a different file to be DRY
// json deal struct for easy AWS upload; ID is md5 hash of location ID + deal
type deal struct {
	ID       string   `json:"id"`
	Location string   `json:"location"`
	Info     string   `json:"info"`
	Days     []string `json:"days"`
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

func getDeals() []deal {
	var deals []deal

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if os.Getenv("LAMBDA_ENV") == "TEST" {
		// Use local
		sess, err = session.NewSession(
			&aws.Config{
				Region:   aws.String(region),
				Endpoint: aws.String(localEndpoint),
			})
	}
	if err != nil {
		log.Print("Error getting session...", err)
		log.Print(err)
		// this might not be good inside lambda function???
		os.Exit(2)
	}

	svc := dynamodb.New(sess)
	params := buildParams()

	result, err := svc.Scan(params)
	if err != nil {
		log.Print("Error scanning...", err)
		os.Exit(2)
	}

	for _, i := range result.Items {
		d := deal{}

		err := dynamodbattribute.UnmarshalMap(i, &d)

		if err != nil {
			log.Print("Error unmarshaling item...", err)
			os.Exit(2)
		}

		deals = append(deals, d)
	}

	return deals
}
func filterByLocation(list []deal, ID string) (out []deal) {
	// use a map to filter out non correct ids
	//https: //stackoverflow.com/questions/32867780/filtering-a-slice-of-structs-based-on-a-different-slice-in-golang#32867892
	// make a map of the one
	f := map[string]bool{ID: true}
	for _, deal := range list {
		if _, ok := f[deal.Location]; ok {
			out = append(out, deal)
		}
	}
	return
}

func filterByDays(list []deal, daysString string) (out []deal) {
	// split the query string
	days := strings.Split(daysString, ",")
	//make a map of the days you want, then check each item and
	// it's not in the map, don't add it
	f := make(map[string]bool, len(days))
	for _, day := range days {
		f[day] = true
	}

	for _, deal := range list {
		// have to check all since each deal can have multiple days
		for _, dealDay := range deal.Days {
			if _, ok := f[dealDay]; ok {
				out = append(out, deal)
				continue
			}
		}
	}
	return
}

// Only location and days for now;
func filterDeals(list []deal, params map[string]string) []deal {
	if _, exists := params["location"]; exists {
		list = filterByLocation(list, params["location"])
	}
	if _, exists := params["days"]; exists {
		list = filterByDays(list, params["days"])
	}
	return list
}

// Handler processes the DynamoDB query response and returns formatted json body
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Fetching deals...")

	if os.Getenv("LAMBDA_ENV") == "TEST" {
		// TODO: error handle if it doesn't exist
		tableName = os.Getenv("TEST_DB")
	}

	deals := getDeals()
	filtered := filterDeals(deals, request.QueryStringParameters)
	log.Print("Filtered deals...")
	marshalled, err := json.Marshal(filtered)
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
