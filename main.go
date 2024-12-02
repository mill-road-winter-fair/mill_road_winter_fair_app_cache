package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

func main() {

	//Parse the argument flags (like the ones in Heroku's Procfile)
	flag.Parse()

	//Get the port from environment variables
	port := os.Getenv("PORT")
	//Define default port value if one is not set
	if port == "" {
		//Use local port
		port = "8080"

		//If we're running locally we also need to load the .env file
		err := godotenv.Load()
		if err != nil {
			glog.Fatal("Error loading .env file")
		}
	}

	// Start the data fetching in a separate goroutine
	glog.Info("Starting fetch of data from Google Sheets API")
	go fetchSheetData()

	//Create default webserver config
	glog.Info("Startign web server")
	webServer := gin.Default()

	//API endpoints to handle shop CRUD operations
	webServer.GET("/listings", GetListingsFromCache)

	//Run the webserver
	ginErr := webServer.Run(":" + port)
	if ginErr != nil {
		glog.Fatalf("Web server initialisation failed: %v", ginErr)
	}
}

func GetListingsFromCache(c *gin.Context) {
	var listingsJson []byte

	//Call the DB function
	glog.Info("Calling getSheetDataFromCache function")
	listingsJson, err := getSheetDataFromCache()

	if err != nil {
		//Return the status code and body from the function
		glog.Error("Returning 500 response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Return successful response with the processed JSON data
	glog.Info("Returning 200 response")
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", listingsJson)
	return
}

var (
	// Shared variable to store the fetched JSON
	sheetData []byte
	// Mutex to synchronize access to `sheetData`
	mu sync.Mutex
)

// getSheetData returns the cached JSON data.
func getSheetDataFromCache() ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	if sheetData == nil {
		glog.Error("sheetData variable is empty")
		return nil, fmt.Errorf("No listings available yet")
	}
	return sheetData, nil
}

// fetchSheetData fetches data from the Google Sheets API at regular intervals.
func fetchSheetData() {
	// Get environment variables
	apiKey := os.Getenv("GOOGLE_SHEETS_API_KEY")
	sheetID := os.Getenv("GOOGLE_SHEET_ID")
	rangeValue := os.Getenv("GOOGLE_SHEET_RANGE")

	if sheetID == "" || apiKey == "" || rangeValue == "" {
		glog.Error("Environment variables GOOGLE_SHEETS_API_KEY, GOOGLE_SHEET_ID, and GOOGLE_SHEET_RANGE must be set.")
		return
	}

	glog.Info("Making HTTP call to Google Sheets API")
	url := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s", sheetID, rangeValue, apiKey)
	
	glog.Info("Beginning one minute delay")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	var lastFetchedData []byte

	glog.Info("Processing respone from Google Sheets API")
	for {
		select {
		case <-ticker.C:
			resp, err := http.Get(url)
			if err != nil {
				glog.Errorf("Error fetching data: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				glog.Errorf("Non-200 response: %d\n", resp.StatusCode)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				glog.Errorf("Error reading response body: %v\n", err)
				continue
			}

			// Check if data has changed
			if string(body) != string(lastFetchedData) {
				glog.Info("Data updated.")
				mu.Lock()
				sheetData = body
				mu.Unlock()
				lastFetchedData = body
			} else {
				glog.Info("No changes in data.")
			}
		}
	}
}
