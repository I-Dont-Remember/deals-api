package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	tablename = "devDeals"
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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := getMapFromJSON(request.Body)

	// if we return error, we get a 502 on the API no matter what
	// so return nil and then HTTP errors
	secretKey := os.Getenv("secret_key")
	secret := os.Getenv("secret")
	if secretKey == "" || secret == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Server doesn't have necessary credentials",
			StatusCode: 500,
		}, nil
	}
	if body[secretKey] != secret {
		return events.APIGatewayProxyResponse{
			Body:       "Forbidden",
			StatusCode: 403,
		}, nil
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Print(err.Error())
	}
	svc := dynamodb.New(sess)

	input, err := formatInput(body, secretKey)
	if err != nil {
		val := strings.Join([]string{"Can't put deal into Dynamo", err.Error()}, " ")
		return events.APIGatewayProxyResponse{
			Body:       val,
			StatusCode: 400,
		}, nil
	}

	_, err = svc.PutItem(input)
	// result appears to just be empty json on success {}
	if err != nil {
		log.Print(err.Error())
	}

	// These are necessary alongside API Gateway CORS enabling
	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
