package models

// Deal is a json helper
type Deal struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	LocationID  string   `json:"location_id"`
	Days        []string `json:"days"`
	AllDay      bool     `json:"all_day"`
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
	Type        []string `json:"type"`
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
	Deals          []string `json:"deals"`
}

// Campus is a json helper
type Campus struct {
	Slug        string   `json:"slug"`
	DisplayName string   `json:"display_name"`
	Locations   []string `json:"locations"`
}
