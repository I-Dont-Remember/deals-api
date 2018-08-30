package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/I-Dont-Remember/deals-api/pkg/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	region    = "us-east-2"
	output    *dynamodb.ScanOutput
	tableName = "newDeals"
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
	filt := expression.Name("id").AttributeExists()
	proj := expression.NamesList(expression.Name("id"),
		expression.Name("days"),
		expression.Name("info"),
		expression.Name("location"))

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

	svc, err := db.CreateConnection(region)
	if err != nil {
		log.Print("Error with creating connection:", err)
	}
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
	deals := getDeals()
	log.Printf("Fetched %d deals from %s", len(deals), tableName)
	filtered := filterDeals(deals, request.QueryStringParameters)
	log.Printf("Filtered deals has %d deals", len(filtered))
	marshalled, err := json.Marshal(filtered)
	if err != nil {
		log.Print("Error marshalling deals...")
		return helpers.ErrResponse("Failed marshalling deals", err, http.StatusInternalServerError)
	}

	return helpers.Response(string(marshalled), http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
