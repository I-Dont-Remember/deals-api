package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// CreateConnnection to DynamoDB
func CreateConnection(region string) (*dynamodb.DynamoDB, error) {
	localEndpoint := "http://localhost:4569/"
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if os.Getenv("API_ENV") == "local" {
		sess, err = session.NewSession(
			&aws.Config{
				Region:   aws.String(region),
				Endpoint: aws.String(localEndpoint),
			})
	}

	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}
