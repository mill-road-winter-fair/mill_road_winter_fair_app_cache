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

var (
	// Shared variables to store the Listings JSON for production and development endpoints
	sheetData []byte
	devSheetData []byte

	// Shared variables to store the Notifications JSON for production and development endpoints
	notificationsSheetData []byte
	devNotificationsSheetData []byte

	// Mutex to synchronize access to listing caches
	mu sync.Mutex

	// ourApiKey holds the expected API key loaded from environment variables.
	ourApiKey string
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

	// Validate the API key
	ValidateApiKey()

	// Start the data fetching in separate goroutines
	glog.Info("Starting fetch of data from Google Sheets API")
	go fetchSheetData()
	go fetchDevSheetData()
	go fetchNotificationsSheetData()
	go fetchDevNotificationsSheetData()

	// Create default webserver config
	glog.Info("Starting web server")
	webServer := gin.Default()

	// API endpoints to handle shop CRUD operations
	webServer.GET("/listings", ListingsEndpoint)
	webServer.GET("/dev-listings", DevListingsEndpoint)
	webServer.GET("/notifications", NotificationsEndpoint)
	webServer.GET("/dev-notifications", DevNotificationsEndpoint)

	// Run the webserver
	ginErr := webServer.Run(":" + port)
	if ginErr != nil {
		glog.Fatalf("Web server initialisation failed: %v", ginErr)
	}
}

func ValidateApiKey() {
	// Get api key environment variable
	ourApiKey = os.Getenv("OUR_API_KEY")
	if ourApiKey == "" {
		glog.Fatal("OUR_API_KEY environment variable is not set")
	}
}

func ListingsEndpoint(c *gin.Context) {
	validateKeyAndReturnListings(c, GetListingsFromCache)
}

func DevListingsEndpoint(c *gin.Context) {
	validateKeyAndReturnListings(c, GetDevListingsFromCache)
}

func NotificationsEndpoint(c *gin.Context) {
	if !requireAPIKey(c) {
		return
	}
	returnCachedListings(c, func() ([]byte, error) {
		return getNotificationsFromCache(&notificationsSheetData)
	})
}

func DevNotificationsEndpoint(c *gin.Context) {
	if !requireAPIKey(c) {
		return
	}
	returnCachedListings(c, func() ([]byte, error) {
		return getNotificationsFromCache(&devNotificationsSheetData)
	})
}

// requireAPIKey uses the same key and header/query precedence as listings,
// but rejects unauthorised requests before accessing the notification cache.
func requireAPIKey(c *gin.Context) bool {
	key := c.GetHeader("X-API-Key")
	if key == "" {
		key = c.Query("key")
	}
	if key == "" || ourApiKey == "" || key != ourApiKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid API key"})
		return false
	}
	return true
}

func getNotificationsFromCache(cache *[]byte) ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	if *cache == nil {
		return nil, fmt.Errorf("No notifications available yet")
	}
	return *cache, nil
}

func validateKeyAndReturnListings(c *gin.Context, endpointFunc gin.HandlerFunc) {
	key := c.GetHeader("X-API-Key")
	if key == "" {
		key = c.Query("key")
	}
	if key == "" {
		glog.Warning("Missing key parameter, returning listings anyway")
		// The first step is to return the listings from the cache, even if the key is missing or invalid. This way, users can still access the data without providing a key, but we will be informed about the missing or invalid key in the logs.
		// Once we have the logs, we can decide whether to enforce the key requirement in the future. To begin with, we will allow access to the listings even if the key is missing or invalid, but we will log a warning message in both cases.
		endpointFunc(c)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "missing key parameter"})
		return
	}
	if key != ourApiKey {
		glog.Warning("Invalid key provided, returning listings anyway")
		// The first step is to return the listings from the cache, even if the key is missing or invalid. This way, users can still access the data without providing a key, but we will be informed about the missing or invalid key in the logs.
		// Once we have the logs, we can decide whether to enforce the key requirement in the future. To begin with, we will allow access to the listings even if the key is missing or invalid, but we will log a warning message in both cases.
		endpointFunc(c)
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid key"})
		return
	}
	glog.Info("Valid key provided, returning listings")
	endpointFunc(c)
}

func GetListingsFromCache(c *gin.Context) {
	returnCachedListings(c, getSheetDataFromCache)
}

func GetDevListingsFromCache(c *gin.Context) {
	returnCachedListings(c, getDevSheetDataFromCache)
}

func returnCachedListings(c *gin.Context, cacheFunc func() ([]byte, error)) {
	// Call the DB function
	glog.Info("Calling cache fetch function")
	listingsJSON, err := cacheFunc()

	if err != nil {
		// Return the status code and body from the function
		glog.Error("Returning 500 response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return successful response with the processed JSON data
	glog.Info("Returning 200 response")
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Data(http.StatusOK, "application/json", listingsJSON)
}

// getSheetDataFromCache returns the cached JSON data for /listings.
func getSheetDataFromCache() ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	if sheetData == nil {
		glog.Error("sheetData variable is empty")
		return nil, fmt.Errorf("No listings available yet")
	}
	return sheetData, nil
}

// getDevSheetDataFromCache returns the cached JSON data for /dev-listings.
func getDevSheetDataFromCache() ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	if devSheetData == nil {
		glog.Error("devSheetData variable is empty")
		return nil, fmt.Errorf("No dev listings available yet")
	}
	return devSheetData, nil
}

// fetchSheetData fetches data from the Google Sheets API at regular intervals.
func fetchSheetData() {
	fetchSheetDataFromEnv("GOOGLE_SHEET_ID", "GOOGLE_SHEET_RANGE", &sheetData)
}

// fetchDevSheetData fetches dev data from a separate Google Sheet at regular intervals.
func fetchDevSheetData() {
	fetchSheetDataFromEnv("DEV_GOOGLE_SHEET_ID", "DEV_GOOGLE_SHEET_RANGE", &devSheetData)
}

func fetchNotificationsSheetData() {
	fetchSheetDataFromEnv("NOTIFICATIONS_GOOGLE_SHEET_ID", "NOTIFICATIONS_GOOGLE_SHEET_RANGE", &notificationsSheetData)
}

func fetchDevNotificationsSheetData() {
	fetchSheetDataFromEnv("DEV_NOTIFICATIONS_GOOGLE_SHEET_ID", "DEV_NOTIFICATIONS_GOOGLE_SHEET_RANGE", &devNotificationsSheetData)
}

func fetchSheetDataFromEnv(sheetIDEnvName, rangeEnvName string, destination *[]byte) {
	// Get environment variables
	apiKey := os.Getenv("GOOGLE_SHEETS_API_KEY")
	sheetID := os.Getenv(sheetIDEnvName)
	rangeValue := os.Getenv(rangeEnvName)

	if sheetID == "" || apiKey == "" || rangeValue == "" {
		glog.Errorf("Environment variables GOOGLE_SHEETS_API_KEY, %s, and %s must be set.", sheetIDEnvName, rangeEnvName)
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
					*destination = body
					mu.Unlock()
					lastFetchedData = body
				} else {
					glog.Info("No changes in data.")
				}
			}()
		}
	}
}
