package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/I-Dont-Remember/deals-api/pkg/db"
	"github.com/I-Dont-Remember/deals-api/pkg/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	tablename = "devLocations"
	region    = "us-east-2"
)

func getMapFromJSON(jsonString string) map[string]string {
	var data map[string]string
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		log.Println("Bad!", err)
	}
	return data
}

func formatInput(dict map[string]string, key string) (*dynamodb.PutItemInput, error) {
	// remove our secret key to prevent it from being displayed anywhere
	delete(dict, key)

	item := make(map[string]*dynamodb.AttributeValue)
	for k, v := range dict {
		item[k] = &dynamodb.AttributeValue{S: aws.String(v)}
	}
	return &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tablename),
	}, nil
}

// Handler is the lambda handler
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := getMapFromJSON(request.Body)

	// if we return error, we get a 502 on the API no matter what
	// so return nil and then HTTP errors
	secretKey := os.Getenv("secret_key")
	secret := os.Getenv("secret")
	if secretKey == "" || secret == "" {
		log.Print("Server doesn't have necessary credentials")
		return helpers.ErrResponse("", nil, http.StatusInternalServerError)

	}
	if body[secretKey] != secret {
		return helpers.ErrResponse("", nil, http.StatusForbidden)
	}

	svc, err := db.CreateConnection(region)
	if err != nil {
		log.Print(err.Error())
	}

	input, err := formatInput(body, secretKey)
	if err != nil {
		return helpers.ErrResponse("", err, http.StatusBadRequest)
	}

	_, err = svc.PutItem(input)
	// result appears to just be empty json on success {}
	if err != nil {
		log.Print(err.Error())
		return helpers.ErrResponse("Can't put location into Dynamo", err, http.StatusInternalServerError)
	}

	return helpers.Response("", http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
