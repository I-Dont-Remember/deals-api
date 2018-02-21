// Loads csv files of Deals into the given table name while providing them
// a unique ksuid value for an Id.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/segmentio/ksuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

// Deal is a rough analogue to the DynamoDB schema; json for easy adding to table
type Deal struct {
	Id       string `json:"Id"`
	Day      string `json:"Day"`
	Location string `json:"Location"`
	Deal     string `json:"Deal"`
}

func validFile(filepath string) bool {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist...\n", filepath)
		} else {
			fmt.Printf("%s exists but some other error occurred...\n", filepath)
		}
		return false
	}

	return true
}

// Return slice of Deals, which get marshalled by AWS function
func getDeals(filepath string) []Deal {
	f, err := os.Open(filepath)
	check(err)
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	check(err)

	deals := make([]Deal, len(lines))
	for i, line := range lines {
		deals[i] = Deal{ksuid.New().String(), line[0], line[1], line[2]}
	}

	return deals
}

func main() {
	var tableName, filepath string
	flag.StringVar(&tableName, "table-name", "", "Name of the table to add items")
	flag.StringVar(&filepath, "path", "", "Path to the csv file of deals to add")
	flag.Parse()

	if tableName == "" || filepath == "" {
		fmt.Println("Need table name and file path...")
		os.Exit(1)
	}

	if !validFile(filepath) {
		os.Exit(1)
	}

	fmt.Println("Initializing AWS session...")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")})
	check(err)

	svc := dynamodb.New(sess)
	deals := getDeals(filepath)

	for _, deal := range deals {
		marshalledDeal, err := dynamodbattribute.MarshalMap(deal)
		check(err)

		input := &dynamodb.PutItemInput{
			Item:      marshalledDeal,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		check(err)

		fmt.Printf("Added %s: %s@%s -> %s\n", deal.Id, deal.Day, deal.Location, deal.Deal)
	}
}
