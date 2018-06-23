package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	LOCATION_TABLE = "devLocations"
	DEAL_TABLE     = "devDeals"
)

// json deal struct for easy AWS upload; ID is md5 hash of location ID + deal
type deal struct {
	ID       string   `json:"id"`
	Location string   `json:"location"`
	Info     string   `json:"info"`
	Days     []string `json:"days"`
}

// json location struct for easy AWS upload; ID is Name + Address
// Address: city,state 5dgitzipcode
type location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// What will we do about locations on multiple campuses?
	Campus  string `json:"campus"`
	Address string `json:"address"`
	//Deals   []deal
}

// struct for toml deal info
type dealInfo struct {
	Info string
	Days []string
}

// struct for location toml file layout
type locationInfo struct {
	Name    string
	Campus  string
	Address string
	Deal    []dealInfo
}

func checkErr(e error) {
	if e != nil {
		fmt.Println("Error: ", e)
		os.Exit(1)
	}
}

func getSVC(notLocal bool) *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	localSess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String("http://localhost:4569/")})

	if notLocal {
		return dynamodb.New(sess)
	}
	return dynamodb.New(localSess)
}

func createLocationID(name, address string) string {
	data := []byte(strings.Join([]string{name, address}, ""))
	hexHash := md5.Sum(data)
	return hex.EncodeToString((hexHash[:]))
}

func createDealID(lid, deal string) string {
	data := []byte(strings.Join([]string{lid, deal}, ""))
	hexHash := md5.Sum(data)
	return hex.EncodeToString(hexHash[:])
}

// uses the dealInfo and location struct to pull out toml contents
func decodeFile(filePath string) (locationInfo, error) {
	var curr locationInfo

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return curr, err
	}

	ret, err := toml.DecodeFile(filePath, &curr)
	if err != nil {
		return curr, err
	}

	if len(ret.Undecoded()) != 0 {
		fmt.Println("[!] unable to decode these correctly", ret.Undecoded())
	}

	return curr, nil
}

// Loop through directory of toml files and read their deal contents into a slice
func getLocationsFromDir(dirPath string) ([]locationInfo, error) {
	var locations []locationInfo

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, err
	}

	// Change to directory since ioutil returns just FileInfo with no path
	os.Chdir(dirPath)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".toml" {

			fmt.Println(file.Name())
			place, err := decodeFile(file.Name())
			if err != nil {
				return nil, err
			}
			locations = append(locations, place)
		}
	}

	return locations, nil
}

func uploadDeal(svc *dynamodb.DynamoDB, d deal) error {
	marshalledDeal, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      marshalledDeal,
		TableName: aws.String(DEAL_TABLE),
	}
	_, err = svc.PutItem(input)
	// just show err and keep going
	if err != nil {
		fmt.Println("[!] Error uploading deal:", err)
	} else {
		fmt.Println("[*] Uploaded deal:", d)
	}

	return nil
}

func uploadLocation(svc *dynamodb.DynamoDB, l location) error {
	marshalledLocation, err := dynamodbattribute.MarshalMap(l)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      marshalledLocation,
		TableName: aws.String(LOCATION_TABLE),
	}
	_, err = svc.PutItem(input)
	// Just show it so we know and keep going
	if err != nil {
		fmt.Println("[!] Error uploading location:", err)
	} else {
		fmt.Println("[*] Uploaded location:", l)

	}

	return nil
}

func uploadItems(locations []location, deals []deal, notLocal bool) error {
	svc := getSVC(notLocal)

	input := &dynamodb.ListTablesInput{}

	result, err := svc.ListTables(input)
	checkErr(err)

	fmt.Println("Existing tables:", result)
	fmt.Printf("Trying to upload %d locations...\n", len(locations))
	for _, l := range locations {
		err := uploadLocation(svc, l)
		checkErr(err)
	}
	fmt.Printf("Trying to upload %d deals...\n", len(deals))
	for _, d := range deals {
		err = uploadDeal(svc, d)
		checkErr(err)
	}
	return nil

}

// TODO: see if we should add error handling to struct creation
func jsonifyLocations(data []locationInfo) ([]location, []deal) {
	// convert all the locationInfo into json structs
	locations := make([]location, len(data))
	deals := []deal{}
	for i, l := range data {
		lID := createLocationID(l.Name, l.Address)

		// convert this specific locations dealInfo into json struct
		// TODO: this is not an efficient way of doing this,
		// but currently not worth time to fix
		for _, d := range l.Deal {
			dID := createDealID(lID, d.Info)
			deals = append(deals, deal{dID, lID, d.Info, d.Days})
		}

		locations[i] = location{lID, l.Name, l.Campus, l.Address}
	}

	return locations, deals
}

func main() {
	var notLocal bool
	var notDev bool
	var path string
	flag.StringVar(&path, "d", "", "/path/to/tomlfile/directory/, required")
	flag.BoolVar(&notLocal, "not-local", false, "Flag to non-local DynamoDB")
	flag.BoolVar(&notDev, "not-dev", false, "Flag to not use the dev tables but the Prod")
	flag.Parse()

	if notDev {
		LOCATION_TABLE = "Locations"
		DEAL_TABLE = "Deals"
	}

	if path == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// loop through directory of toml files to gather all locations & deals
	allLocationInfo, err := getLocationsFromDir(path)
	checkErr(err)

	locations, deals := jsonifyLocations(allLocationInfo)

	// upload all the deals, have to do locations first to make sure they exist and we aren't referencing empty ID's
	err = uploadItems(locations, deals, notLocal)
	checkErr(err)
}
