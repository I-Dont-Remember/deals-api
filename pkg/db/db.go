// Package db is a wrapper to coalesce common database intricacies and let the callee
// ignore how it's done
package db

import (
	"log"
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
	GetCampus(slug string) (models.Campus, error)
	CreateLocation(models.Location) error
	RemoveLocation(id string) error
	GetLocations() ([]models.Location, error)
	GetLocation(id string) (models.Location, error)
	GetDeals() ([]models.Deal, error)
	RemoveDeal(id string) error
	CreateDeal(models.Deal) error
}

// TODO: the get functions should be switched to using pointers
// TODO: pull out the Dynamo getItem into it's own function since it's used so much

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

// GetLocations fetches all locations
func (db Dynamo) GetLocations() ([]models.Location, error) {
	locations := []models.Location{}

	si := &dynamodb.ScanInput{
		TableName: aws.String(locationTable),
	}
	result, err := db.conn.Scan(si)
	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		l := models.Location{}
		err := dynamodbattribute.UnmarshalMap(i, &l)

		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

// GetLocations - mocked
func (m Mock) GetLocations() ([]models.Location, error) {
	return []models.Location{models.Location{ID: "from the mock"}}, nil
}

// GetLocation fetches a single location based on id
func (db Dynamo) GetLocation(id string) (models.Location, error) {
	l := &models.Location{}

	gi := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(locationTable),
	}
	result, err := db.conn.GetItem(gi)
	if err != nil {
		return *l, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, l)
	if err != nil {
		return *l, err
	}

	return *l, nil
}

// GetLocation - mocked
func (m Mock) GetLocation(id string) (models.Location, error) {
	return models.Location{}, nil
}

// CreateLocation makes a new Location
func (db Dynamo) CreateLocation(l models.Location) error {
	av, err := dynamodbattribute.MarshalMap(l)
	if err != nil {
		return err
	}

	pi := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(locationTable),
	}

	_, err = db.conn.PutItem(pi)
	return err
}

// CreateLocation - mocked
func (m Mock) CreateLocation(models.Location) error {
	return nil
}

// RemoveLocation removes the location
func (db Dynamo) RemoveLocation(id string) error {
	di := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(locationTable),
	}
	log.Print(di)
	// TODO: this returns the same val,err when using a non-existent key, decide if we should throw error on invalid key
	val, err := db.conn.DeleteItem(di)
	log.Print(val)
	log.Print(err)
	return err
}

// RemoveLocation - mocked
func (m Mock) RemoveLocation(id string) error {
	return nil
}

// GetDeals fetches all deals
func (db Dynamo) GetDeals() ([]models.Deal, error) {
	deals := []models.Deal{}

	si := &dynamodb.ScanInput{
		TableName: aws.String(dealTable),
	}
	result, err := db.conn.Scan(si)
	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		d := models.Deal{}
		err := dynamodbattribute.UnmarshalMap(i, &d)

		if err != nil {
			return nil, err
		}
		deals = append(deals, d)
	}
	return deals, nil
}

// GetDeals - mocked
func (m Mock) GetDeals() ([]models.Deal, error) {
	return []models.Deal{models.Deal{ID: "from the mock"}}, nil
}

// CreateDeal makes a new Deal
func (db Dynamo) CreateDeal(d models.Deal) error {
	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return err
	}

	pi := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dealTable),
	}

	_, err = db.conn.PutItem(pi)
	return err
}

// CreateDeal - mocked
func (m Mock) CreateDeal(d models.Deal) error {
	return nil
}

// RemoveDeal removes the location
func (db Dynamo) RemoveDeal(id string) error {
	di := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(dealTable),
	}
	_, err := db.conn.DeleteItem(di)
	return err
}

// RemoveDeal - mocked
func (m Mock) RemoveDeal(id string) error {
	return nil
}
