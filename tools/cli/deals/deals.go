// Package deals will act as a simple SDK for the functions we need until we find a better solution
package deals

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/I-Dont-Remember/deals-api/pkg/models"
)

// Client is an API client
type Client struct {
	basePath string
}

// New creates a new API client
func New(basePath string) *Client {
	return &Client{
		basePath: basePath,
	}
}

func (c *Client) post(path string, data map[string]interface{}) ([]byte, error) {
	url := c.basePath + path

	str, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println("Making API call to " + url)
	body := bytes.NewBuffer(str)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Dot-Auth", "local")
	response, postErr := (&http.Client{}).Do(request)
	if postErr != nil {
		return []byte{}, errors.New("Failed during post - " + postErr.Error())
	}

	respBody, _ := ioutil.ReadAll(response.Body)

	return respBody, nil
}

// CreateCampus creates a campus with API
func (c *Client) CreateCampus(slug, name string) (models.Campus, error) {
	data := map[string]interface{}{
		"slug":         slug,
		"display_name": name,
	}
	campus := models.Campus{}

	// TODO: these are the exact same process for all post calls expecting interface returns,
	// this could be done by figuring out how to just pass it route,
	// data, and the pointer to be filled by json
	body, err := c.post("/campuses", data)
	if err != nil {
		return campus, err
	}

	if err := json.Unmarshal(body, &campus); err != nil {
		fmt.Println(string(body))
		return campus, err
	}

	return campus, nil
}

// CreateLocation creates a location with API
func (c *Client) CreateLocation(slug, name, address, lat, long, imageLink, phoneNumber, website, yelpLink string) (models.Location, error) {
	path := "/campuses/" + slug + "/locations"
	data := map[string]interface{}{
		"name":            name,
		"display_address": address,
		"latitude":        lat,
		"longitude":       long,
		"image_link":      imageLink,
		"phone_number":    phoneNumber,
		"website":         website,
		"yelp_link":       yelpLink,
	}
	location := models.Location{}

	body, err := c.post(path, data)
	if err != nil {
		return location, err
	}

	if err := json.Unmarshal(body, &location); err != nil {
		fmt.Println(string(body))
		return location, err
	}

	return location, nil
}

// CreateDeal creates a deal with API
func (c *Client) CreateDeal(locationID, desc, start, end string, days, types []string, allDay bool) (models.Deal, error) {
	path := "/locations/" + locationID + "/deals"
	data := map[string]interface{}{
		"description": desc,
		"start_time":  start,
		"end_time":    end,
		"all_day":     allDay,
		"days":        days,
		"types":       types,
	}
	deal := models.Deal{}

	body, err := c.post(path, data)
	if err != nil {
		return deal, err
	}

	if err := json.Unmarshal(body, &deal); err != nil {
		fmt.Println(string(body))
		return deal, err
	}
	return deal, nil
}
