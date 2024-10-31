package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/listings"
	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/shared"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
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
		message, err := addTestingDataToLocalDb(db)
		if err != nil {
			glog.Fatalf("Error adding test data: %v", err)
		} else {
			fmt.Println(message)
		}
	} else {
		fmt.Println("Please specify a valid operation using the flags.")
	}
}

// Function to add test data to the local database
func addTestingDataToLocalDb(db *sql.DB) (details string, err error) {

	//Hardcoded slice of ShopData structs
	testListings := []shared.ListingData{
		{Name: "glazedandconfused", DisplayName: "Glazed and Confused", PrimaryType: "Vendor", SecondaryType: "Food", TertiaryType: "Doughnuts", Email: "admin@glazedandconfued.com", Website: "https://www.glazedandconfused.com", Phone: "01223 111111", PlusCode: "9F4254XQ+VG", StartTime: "09:00", EndTime: "15:00"},
		{Name: "sushisquad", DisplayName: "Sushi Squad", PrimaryType: "Vendor", SecondaryType: "Food", TertiaryType: "Sushi", Email: "admin@sushisquad.com", Website: "https://www.sushisquad.com", Phone: "01223 222222", PlusCode: "9F4254XQ+XMG", StartTime: "09:00", EndTime: "15:00"},
		{Name: "familyjewels", DisplayName: "Family Jewels", PrimaryType: "Vendor", SecondaryType: "Retail", TertiaryType: "Jewellery", Email: "admin@familyjewels.com", Website: "https://www.familyjewels.com", Phone: "01223 333333", PlusCode: "9F42642M+J8", StartTime: "09:00", EndTime: "16:00"},
		{Name: "rollingstones", DisplayName: "The Rolling Stones", PrimaryType: "Performer", SecondaryType: "Musician", TertiaryType: "Rock", Email: "admin@rollingstones.com", Website: "https://rollingstones.com/", Phone: "01223 444444", PlusCode: "9F4254XW+23", StartTime: "09:30", EndTime: "10:30"},
		{Name: "dukeellington", DisplayName: "Duke Ellington", PrimaryType: "Performer", SecondaryType: "Musician", TertiaryType: "Jazz", Email: "duke@ellington.com", Website: "https://en.wikipedia.org/wiki/Duke_Ellington", Phone: "01223 555555", PlusCode: "9F4254XW+23", StartTime: "10:30", EndTime: "11:30"},
		{Name: "muddywaters", DisplayName: "Muddy Waters", PrimaryType: "Performer", SecondaryType: "Musician", TertiaryType: "Blues", Email: "muddy@waters.com", Website: "https://en.wikipedia.org/wiki/Muddy_Waters", Phone: "01223 666666", PlusCode: "9F4254XW+23", StartTime: "11:30", EndTime: "12:30"},
		{Name: "knittingcircle", DisplayName: "Knitting Circle", PrimaryType: "Event", SecondaryType: "Craft", TertiaryType: "Knitting", Email: "", Website: "https://www.theknittingnetwork.co.uk/", Phone: "01223 777777", PlusCode: "9F4254XQ+R6", StartTime: "09:30", EndTime: "10:30"},
		{Name: "standupcomedy", DisplayName: "Stand Up Comedy", PrimaryType: "Event", SecondaryType: "Performance", TertiaryType: "Comedy", Email: "", Website: "", Phone: "01223 888888", PlusCode: "9F4254XR+CQ", StartTime: "13:30", EndTime: "15:30"},
		{Name: "publictoilets", DisplayName: "Public Toilets", PrimaryType: "Service", SecondaryType: "Toilets", TertiaryType: "Accessible", Email: "public@toilets.com", Website: "", Phone: "01223 999999", PlusCode: "9F4254XQ+RCX", StartTime: "09:00", EndTime: "17:00"},
	}

	//Iterate through the testShops slice and insert each shop into the DB
	for _, listing := range testListings {
		_, details, err := listings.CreateListingInDb(db, listing.Name, listing.DisplayName, listing.PrimaryType, listing.SecondaryType, listing.TertiaryType, listing.Email, listing.Website, listing.Phone, listing.PlusCode, listing.StartTime, listing.EndTime)
		if err != nil {
			glog.Errorf("Error adding testListings data: %v", details)
			return details, err
		}
	}

	return "Added test data to local DB successfully", nil
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
