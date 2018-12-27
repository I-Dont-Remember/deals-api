package db

import (
	"os"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DB is a wrapper to allow easy mocking & swapping of persistent storage
type DB interface {
	GetCampuses() []models.Campus
	GetCampus(string) models.Campus
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
func (db Dynamo) GetCampuses() []models.Campus {
	return []models.Campus{models.Campus{Slug: "slug campuses"}}
}

// GetCampuses - mocked
func (m Mock) GetCampuses() []models.Campus {
	return []models.Campus{models.Campus{Slug: "from the mock"}}
}

// GetCampus fetches a single campus based on slug
func (db Dynamo) GetCampus(slug string) models.Campus {
	return models.Campus{Slug: "slug"}
}

// GetCampus - mocked
func (m Mock) GetCampus(slug string) models.Campus {
	return models.Campus{}
}
