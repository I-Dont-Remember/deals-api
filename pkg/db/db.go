// db is a wrapper to coalesce common database intricacies and let the functions
// ignore how it's done
package db

import (
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	campusTable   = "Campuses"
	dealTable     = "Deals"
	locationTable = "Locations"
)

// DB is a wrapper to allow easy mocking & swapping of persistent storage
type DB interface {
	GetCampuses() ([]models.Campus, error)
	GetCampus(string) (models.Campus, error)
}

// Dynamo implements DB
type Dynamo struct {
	conn *dynamodb.DynamoDB
}

// TODO: include data structures in these functions
// that can be edited by unit tests without having
// to edit this file each time

// Mock mocks DB
type Mock struct{}

// Connect returns a Dynamo connection; local, remote, or a Mock
func Connect() (DB, error) {
	region := "us-east-2"
	localEndpoint := "http://localhost:4569/"
	env := os.Getenv("API_ENV")

	if env == "test" {
		return Mock{}, nil
	}

	if env != "local" && env != "prod" {
		// panic since we make important decisions based on env type
		// and don't want any uncertainties
		panic("Unknown API_ENV choice '" + env + "'")
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if env == "local" {
		sess, err = session.NewSession(
			&aws.Config{
				Region:   aws.String(region),
				Endpoint: aws.String(localEndpoint),
			})
	}

	if err != nil {
		return nil, err
	}
	return &Dynamo{conn: dynamodb.New(sess)}, nil
}

// GetCampuses fetches all campuses
func (db Dynamo) GetCampuses() ([]models.Campus, error) {
	campuses := []models.Campus{}

	si := &dynamodb.ScanInput{
		TableName: aws.String(campusTable),
	}
	result, err := db.conn.Scan(si)
	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		c := models.Campus{}
		err := dynamodbattribute.UnmarshalMap(i, &c)

		if err != nil {
			return nil, err
		}
		campuses = append(campuses, c)
	}
	return campuses, nil
}

// GetCampuses - mocked
func (m Mock) GetCampuses() ([]models.Campus, error) {
	return []models.Campus{models.Campus{Slug: "from the mock"}}, nil
}

// GetCampus fetches a single campus based on slug
func (db Dynamo) GetCampus(slug string) (models.Campus, error) {
	c := &models.Campus{}

	gi := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"slug": {
				S: aws.String(slug),
			},
		},
		TableName: aws.String(campusTable),
	}
	result, err := db.conn.GetItem(gi)
	if err != nil {
		return *c, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, c)
	if err != nil {
		return *c, err
	}

	return *c, nil
}

// GetCampus - mocked
func (m Mock) GetCampus(slug string) (models.Campus, error) {
	return models.Campus{}, nil
}
