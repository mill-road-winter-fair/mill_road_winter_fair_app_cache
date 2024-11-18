package shared

import "database/sql"

// Define variable for database (which will be injected from main package)
var Database *sql.DB

type ListingData struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	DisplayName   string `json:"displayName"`
	PrimaryType   string `json:"primaryType"`
	SecondaryType string `json:"secondaryType"`
	TertiaryType  string `json:"tertiaryType"`
	Email         string `json:"email"`
	Website       string `json:"website"`
	Phone         string `json:"phone"`
	LatLng        string `json:"latLng"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
}
