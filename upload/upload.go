package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
	Deals   []deal
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

func createLocationID() string {
	return "0"
}

func createDealID() string {
	return "0"
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

func uploadDeals() {
	fmt.Println("here's where ya upload")
}

// TODO: see if we should add error handling to struct creation
func jsonifyLocations(data []locationInfo) []location {
	// convert all the locationInfo into json structs
	lID := createLocationID()
	dID := createDealID()

	locations := make([]location, len(data))
	for i, l := range data {
		// convert this specific locations dealInfo into json struct
		deals := make([]deal, len(l.Deal))
		for j, d := range l.Deal {
			deals[j] = deal{dID, lID, d.Info, d.Days}
		}
		locations[i] = location{lID, l.Name, l.Campus, l.Address, deals}
	}

	return locations
}

func main() {
	dirPath := "../resources/location-files/"

	// loop through directory of toml files to gather all locations & deals
	allLocationInfo, err := getLocationsFromDir(dirPath)
	checkErr(err)

	fmt.Println(allLocationInfo)

	locations := jsonifyLocations(allLocationInfo)

	fmt.Println(locations)

	// upload all the deals, have to do locations first to make sure they exist and we aren't referencing empty ID's

}
