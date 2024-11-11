package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/listings"
	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/shared"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {

	//Define flags
	wipeDb := flag.Bool("wipe", false, "Wipe the local database")
	addData := flag.Bool("add-data", false, "Add test data to the local database")

	//Parse flags
	flag.Parse()

	//Define connection string
	dsn := "postgres://mrwfadmin:petersfield@localhost/mrwf_db?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	//Check which operation to perform based on the flags
	if *wipeDb {
		message, err := wipeLocalDb(db)
		if err != nil {
			glog.Fatalf("Error wiping DB: %v", err)
		} else {
			fmt.Println(message)
		}
	} else if *addData {

		//Load the .env file
		err := godotenv.Load()
		if err != nil {
			glog.Fatal("Error loading .env file")
		}

		//Get the Google Sheets API key
		googleSheetsAPIKey := os.Getenv("GOOGLE_SHEETS_API_KEY")
		if googleSheetsAPIKey == "" {
			glog.Fatal("GOOGLE_SHEETS_API_KEY is not set in .env file")
		}

		message, err := addTestingDataToLocalDb(db, "1-Dk_K8tvDJ4C9vSx0OJSEYhvhGrt6IEkabVRP83n0OM", "A2:K100", googleSheetsAPIKey)

		if err != nil {
			glog.Fatalf("Error adding test data: %v", err)
		} else {
			fmt.Println(message)
		}
	} else {
		fmt.Println("Please specify a valid operation using the flags.")
	}
}

//Function to add the Google Sheets test data to the local database
func addTestingDataToLocalDb(db *sql.DB, spreadsheetId, rangeName, apiKey string) (string, error) {
	// Fetch data from Google Sheets
	listingsFromSheet, err := fetchListingsFromSheet(spreadsheetId, rangeName, apiKey)
	if err != nil {
		glog.Errorf("Error fetching listings from Google Sheets: %v", err)
		return "Failed to add data from Google Sheets", err
	}

	//Insert each listing into the DB
	for _, listing := range listingsFromSheet {
		_, details, err := listings.CreateListingInDb(db, listing.Name, listing.DisplayName, listing.PrimaryType, listing.SecondaryType, listing.TertiaryType, listing.Email, listing.Website, listing.Phone, listing.PlusCode, listing.StartTime, listing.EndTime)
		if err != nil {
			glog.Errorf("Error adding listing to DB: %v", details)
			return details, err
		}
	}

	return "Added data from Google Sheets to local DB successfully", nil
}

//Function to fetch listings from Google Sheets
func fetchListingsFromSheet(spreadsheetId, rangeName, apiKey string) ([]shared.ListingData, error) {
	ctx := context.Background()

	//Set up Sheets service with API key
	srv, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	//Read values from specified range
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, rangeName).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	var listings []shared.ListingData
	for i, row := range resp.Values {
		if i == 0 {
			// Skip header row
			continue
		}

		if len(row) < 10 {
			glog.Errorf("Skipping row %d: not enough columns", i+1)
			continue
		}

		// Parse each column into ListingData
		listing := shared.ListingData{
			Name:          fmt.Sprintf("%v", row[0]),
			DisplayName:   fmt.Sprintf("%v", row[1]),
			PrimaryType:   fmt.Sprintf("%v", row[2]),
			SecondaryType: fmt.Sprintf("%v", row[3]),
			TertiaryType:  fmt.Sprintf("%v", row[4]),
			Email:         fmt.Sprintf("%v", row[5]),
			Website:       fmt.Sprintf("%v", row[6]),
			Phone:         fmt.Sprintf("%v", row[7]),
			PlusCode:      fmt.Sprintf("%v", row[8]),
			StartTime:     fmt.Sprintf("%v", row[9]),
			EndTime:       fmt.Sprintf("%v", row[10]),
		}

		listings = append(listings, listing)
	}

	return listings, nil
}

// Function to wipe the local database
func wipeLocalDb(db *sql.DB) (details string, err error) {
	//SQL query to wipe all the table
	query := "TRUNCATE TABLE listings RESTART IDENTITY CASCADE;"

	//Execute the query
	_, err = db.Exec(query)
	if err != nil {
		return "Failed to wipe local DB", err
	}

	return "Local DB wiped successfully", nil
}
