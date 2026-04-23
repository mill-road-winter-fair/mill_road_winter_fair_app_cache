package main

import (
	"bytes"
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

	// Parse the argument flags (like the ones in Heroku's Procfile)
	flag.Parse()

	// Get the port from environment variables
	port := os.Getenv("PORT")
	//Define default port value if one is not set
	if port == "" {
		//Use local port
		port = "8080"
	}

	if port == "8080" {
		// If we're running locally we also need to load the .env file
		err := godotenv.Load()
		if err != nil {
			glog.Fatal("Error loading .env file")
		}
	}

	if port != "8080" {
		// If we're not running locally set Gin to Release mode
		gin.SetMode(gin.ReleaseMode)
	}

	// Start the data fetching in a separate goroutine
	glog.Info("Starting fetch of data from Google Sheets API")
	go fetchSheetData()

	// Create default webserver config
	glog.Info("Starting web server")
	webServer := gin.Default()

	// API endpoint to handle listings GET operations
	webServer.GET("/listings", ListingsEndpoint)

	// Run the webserver
	ginErr := webServer.Run(":" + port)
	if ginErr != nil {
		glog.Fatalf("Web server initialisation failed: %v", ginErr)
	}
}

func ListingsEndpoint(c *gin.Context) {
	// Get environment variables
	ourApiKey := os.Getenv("OUR_API_KEY")

	key := c.GetHeader("X-API-Key")
	if key == "" {
		key = c.Query("key")
	}
	if key == "" {
		glog.Warning("Missing key parameter")
		// The first step is to return the listings from the cache, even if the key is missing or invalid. This way, users can still access the data without providing a key, but we will be informed about the missing or invalid key in the logs.
		// Once we have the logs, we can decide whether to enforce the key requirement in the future. To begin with, we will allow access to the listings even if the key is missing or invalid, but we will log a warning message in both cases.
		GetListingsFromCache(c)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "missing key parameter"})
		return
	}
	if key != ourApiKey {
		glog.Warning("Invalid key provided")
		// The first step is to return the listings from the cache, even if the key is missing or invalid. This way, users can still access the data without providing a key, but we will be informed about the missing or invalid key in the logs.
		// Once we have the logs, we can decide whether to enforce the key requirement in the future. To begin with, we will allow access to the listings even if the key is missing or invalid, but we will log a warning message in both cases.
		GetListingsFromCache(c)
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid key"})
		return
	}
	glog.Info("Valid key provided, returning listings")
	GetListingsFromCache(c)
}

func GetListingsFromCache(c *gin.Context) {
	var listingsJson []byte

	// Call the DB function
	glog.Info("Calling getSheetDataFromCache function")
	listingsJson, err := getSheetDataFromCache()

	if err != nil {
		// Return the status code and body from the function
		glog.Error("Returning 500 response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return successful response with the processed JSON data
	glog.Info("Returning 200 response")
	c.Header("Content-Type", "application/json; charset=UTF-8")
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
	googleSheetsApiKey := os.Getenv("GOOGLE_SHEETS_API_KEY")
	googleSheetID := os.Getenv("GOOGLE_SHEET_ID")
	googleSheetRange := os.Getenv("GOOGLE_SHEET_RANGE")

	if googleSheetID == "" || googleSheetsApiKey == "" || googleSheetRange == "" {
		glog.Error("Environment variables GOOGLE_SHEETS_API_KEY, GOOGLE_SHEET_ID, and GOOGLE_SHEET_RANGE must be set.")
		return
	}

	glog.Info("Making HTTP call to Google Sheets API")
	url := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s", googleSheetID, googleSheetRange, googleSheetsApiKey)

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

			func() {
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					glog.Errorf("Non-200 response: %d\n", resp.StatusCode)
					return
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					glog.Errorf("Error reading response body: %v\n", err)
					return
				}

				if !bytes.Equal(body, lastFetchedData) {
					glog.Info("Data updated.")
					mu.Lock()
					sheetData = body
					mu.Unlock()
					lastFetchedData = body
				} else {
					glog.Info("No changes in data.")
				}
			}()
		}
	}
}
