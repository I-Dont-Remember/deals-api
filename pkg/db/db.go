// Package db is a wrapper to coalesce common database intricacies and let the callee
// ignore how it's done
package db

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/I-Dont-Remember/deals-api/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DB is a wrapper to allow easy mocking & swapping of persistent storage
type DB interface {
	InputSearchAnalytics(s models.SearchData) error
	CreateCampus(models.Campus) error
	RemoveCampus(slug string) error
	GetCampuses() ([]models.Campus, error)
	GetCampus(slug string) (models.Campus, error)
	UpdateCampus(models.Campus) (models.Campus, error)
	CreateLocation(models.Location) error
	RemoveLocation(id string) error
	GetLocations() ([]models.Location, error)
	GetLocation(id string) (models.Location, error)
	UpdateLocation(models.Location) (models.Location, error)
	GetDeals() ([]models.Deal, error)
	BatchGetDeals(ids []string) ([]models.Deal, error)
	RemoveDeal(id string) error
	CreateDeal(models.Deal) error
}

// TODO: the get functions should be switched to using pointers
// TODO: pull out the Dynamo getItem into it's own function since it's used so much

// Dynamo implements DB
type Dynamo struct {
	conn           *dynamodb.DynamoDB
	campusTable    string
	dealTable      string
	locationTable  string
	analyticsTable string
}

// Mock mocks DB
type Mock struct {
	CreateCampusFunc         func(models.Campus) error
	RemoveCampusFunc         func(slug string) error
	GetCampusesFunc          func() ([]models.Campus, error)
	GetCampusFunc            func(slug string) (models.Campus, error)
	UpdateCampusFunc         func(models.Campus) (models.Campus, error)
	CreateLocationFunc       func(models.Location) error
	RemoveLocationFunc       func(id string) error
	GetLocationsFunc         func() ([]models.Location, error)
	GetLocationFunc          func(id string) (models.Location, error)
	UpdateLocationFunc       func(models.Location) (models.Location, error)
	GetDealsFunc             func() ([]models.Deal, error)
	BatchGetDealsFunc        func(ids []string) ([]models.Deal, error)
	RemoveDealFunc           func(id string) error
	CreateDealFunc           func(models.Deal) error
	InputSearchAnalyticsFunc func(s models.SearchData) error
}

// Connect returns a Dynamo connection; local or remote
func Connect() (DB, error) {
	region := "us-east-2"
	localEndpoint := "http://localhost:4569/"
	env := os.Getenv("API_ENV")

	if env != "local" && env != "prod" && env != "dev" {
		// panic since we make important decisions based on env type
		// and don't want any uncertainties
		panic("Unknown API_ENV choice '" + env + "'")
	}

	d := &Dynamo{}
	if env == "prod" || env == "local" {
		d.campusTable = "Campuses"
		d.dealTable = "Deals"
		d.locationTable = "Locations"
		d.analyticsTable = "Analytics"
	} else if env == "dev" {
		d.campusTable = "Campuses-dev"
		d.dealTable = "Deals-dev"
		d.locationTable = "Locations-dev"
		d.analyticsTable = "Analytics-dev"
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
	d.conn = dynamodb.New(sess)
	return d, nil
}

// InputSearchAnalytics used for quickly getting search analytics til we have a better way
func (db Dynamo) InputSearchAnalytics(s models.SearchData) error {
	av, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", db)
	pi := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.analyticsTable),
	}

	_, err = db.conn.PutItem(pi)
	return err
}

//InputSearchAnalytics - needed mock because we have to keep interface the same; laymee
func (m Mock) InputSearchAnalytics(s models.SearchData) error {
	return m.InputSearchAnalyticsFunc(s)
}

// CreateCampus makes a new Campus
func (db Dynamo) CreateCampus(c models.Campus) error {
	av, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return err
	}

	pi := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.campusTable),
	}

	_, err = db.conn.PutItem(pi)
	return err
}

// CreateCampus - mocked
func (m Mock) CreateCampus(c models.Campus) error {
	return m.CreateCampusFunc(c)
}

// RemoveCampus removes the campus
func (db Dynamo) RemoveCampus(slug string) error {
	di := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(slug),
			},
		},
		TableName: aws.String(db.campusTable),
	}
	// TODO: this returns the same val,err when using a non-existent key, decide if we should throw error on invalid key
	_, err := db.conn.DeleteItem(di)
	return err
}

// RemoveCampus - mocked
func (m Mock) RemoveCampus(slug string) error {
	return m.RemoveCampusFunc(slug)
}

// GetCampuses fetches all campuses
func (db Dynamo) GetCampuses() ([]models.Campus, error) {
	campuses := []models.Campus{}

	si := &dynamodb.ScanInput{
		TableName: aws.String(db.campusTable),
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
	return m.GetCampusesFunc()
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
		TableName: aws.String(db.campusTable),
	}
	result, err := db.conn.GetItem(gi)
	if err != nil {
		return *c, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, c)
	if err != nil {
		return *c, err
	}

	// // raise error if that campus didn't exist
	// if c.Slug == "" {
	// 	return *c, errors.New("Couldn't find slug -> " + slug)
	// }

	return *c, nil
}

// GetCampus - mocked
func (m Mock) GetCampus(slug string) (models.Campus, error) {
	return m.GetCampusFunc(slug)
}

// UpdateCampus updates the campus's data
func (db Dynamo) UpdateCampus(campus models.Campus) (models.Campus, error) {
	return campus, db.CreateCampus(campus)
}

// UpdateCampus - mocked
func (m Mock) UpdateCampus(campus models.Campus) (models.Campus, error) {
	return m.UpdateCampusFunc(campus)
}

// GetLocations fetches all locations
func (db Dynamo) GetLocations() ([]models.Location, error) {
	locations := []models.Location{}

	si := &dynamodb.ScanInput{
		TableName: aws.String(db.locationTable),
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
	return m.GetLocationsFunc()
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
		TableName: aws.String(db.locationTable),
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
	return m.GetLocationFunc(id)
}

// UpdateLocation updates a location's data
func (db Dynamo) UpdateLocation(location models.Location) (models.Location, error) {
	// TODO: check if there is a better way
	// DynamoDB will "update" an item by replacing the same one, giving the same effect
	return location, db.CreateLocation(location)
}

// UpdateLocation - mocked
func (m Mock) UpdateLocation(location models.Location) (models.Location, error) {
	return m.UpdateLocationFunc(location)
}

// CreateLocation makes a new Location
func (db Dynamo) CreateLocation(l models.Location) error {
	av, err := dynamodbattribute.MarshalMap(l)
	if err != nil {
		return err
	}

	pi := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.locationTable),
	}

	_, err = db.conn.PutItem(pi)
	return err
}

// CreateLocation - mocked
func (m Mock) CreateLocation(l models.Location) error {
	return m.CreateLocationFunc(l)
}

// RemoveLocation removes the location
func (db Dynamo) RemoveLocation(id string) error {
	di := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.locationTable),
	}
	// TODO: this returns the same val,err when using a non-existent key, decide if we should throw error on invalid key
	_, err := db.conn.DeleteItem(di)
	return err
}

// RemoveLocation - mocked
func (m Mock) RemoveLocation(id string) error {
	return m.RemoveLocationFunc(id)
}

// GetDeals fetches all deals
func (db Dynamo) GetDeals() ([]models.Deal, error) {
	deals := []models.Deal{}

	si := &dynamodb.ScanInput{
		TableName: aws.String(db.dealTable),
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
	return m.GetDealsFunc()
}

// BatchGetDeals get deals based on list of ids
func (db Dynamo) BatchGetDeals(ids []string) ([]models.Deal, error) {
	keys := []map[string]*dynamodb.AttributeValue{}
	for _, id := range ids {
		key := map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		}
		keys = append(keys, key)
	}

	requestItems := map[string]*dynamodb.KeysAndAttributes{}
	requestItems[db.dealTable] = &dynamodb.KeysAndAttributes{Keys: keys}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	}
	output, err := db.conn.BatchGetItem(input)
	if err != nil {
		return nil, err
	}

	dynamoDeals, ok := output.Responses[db.dealTable]
	if !ok {
		// didn't have our table in response; is this an error?
		return []models.Deal{}, nil
	}

	var deals []models.Deal

	for _, d := range dynamoDeals {
		deal := models.Deal{}
		err = dynamodbattribute.UnmarshalMap(d, &deal)
		if err != nil {
			return []models.Deal{}, err
		}
		deals = append(deals, deal)
	}

	return deals, nil
}

// BatchGetDeals - mocked
func (m Mock) BatchGetDeals(ids []string) ([]models.Deal, error) {
	return m.BatchGetDealsFunc(ids)
}

// CreateDeal makes a new Deal
func (db Dynamo) CreateDeal(d models.Deal) error {
	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return err
	}

	pi := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.dealTable),
	}

	_, err = db.conn.PutItem(pi)
	return err
}

// CreateDeal - mocked
func (m Mock) CreateDeal(d models.Deal) error {
	return m.CreateDealFunc(d)
}

// RemoveDeal removes the location
func (db Dynamo) RemoveDeal(id string) error {
	di := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.dealTable),
	}
	_, err := db.conn.DeleteItem(di)
	return err
}

// RemoveDeal - mocked
func (m Mock) RemoveDeal(id string) error {
	return m.RemoveDealFunc(id)
}
