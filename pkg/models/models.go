package models

// TODO: omitempty can be used with json, but do we want it?  How do we want to handle non-existent keys?

// Deal is a json helper
// Days options: [Mon, Tue, Wed, Thu, Fri, Sat, Sun]
// Type options: [Drinks, Food, Event]
type Deal struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	LocationID  string   `json:"location_id"`
	Days        []string `json:"days"`
	AllDay      bool     `json:"all_day"`
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
	Types       []string `json:"types"`
}

// Location is a json helper
type Location struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CampusSlug     string   `json:"campus_slug"`
	DisplayAddress string   `json:"display_address"`
	Latitude       string   `json:"latitude"`
	Longitude      string   `json:"longitude"`
	ImageLink      string   `json:"image_link"`
	PhoneNumber    string   `json:"phone_number"`
	Website        string   `json:"website"`
	YelpLink       string   `json:"yelp_link"`
	Hours          []string `json:"hours"`
	Deals          []string `json:"deals"`
}

// ExpandedLocation is a way to trade deal id's for the actual object
type ExpandedLocation struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CampusSlug     string   `json:"campus_slug"`
	DisplayAddress string   `json:"display_address"`
	Latitude       string   `json:"latitude"`
	Longitude      string   `json:"longitude"`
	ImageLink      string   `json:"image_link"`
	PhoneNumber    string   `json:"phone_number"`
	Website        string   `json:"website"`
	YelpLink       string   `json:"yelp_link"`
	Hours          []string `json:"hours"`
	Deals          []Deal   `json:"deals"`
}

// Campus is a json helper
type Campus struct {
	Slug        string   `json:"slug"`
	DisplayName string   `json:"display_name"`
	Locations   []string `json:"locations"`
}
