package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testSheetData = []byte(`[{"name": "Test Listing"}]`)
var testDevSheetData = []byte(`[{"name": "Test Dev Listing"}]`)

// Setup function for tests
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/listings", GetListingsFromCache)
	router.GET("/dev-listings", GetDevListingsFromCache)
	return router
}

// Test GetListingsFromCache when data is available
func TestGetListingsFromCache_Success(t *testing.T) {
	// Set up test data
	mu.Lock()
	sheetData = testSheetData
	mu.Unlock()

	router := setupRouter()

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/listings", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header().Get("Content-Type"))
	assert.JSONEq(t, string(testSheetData), resp.Body.String())
}

// Test GetListingsFromCache when data is not available
func TestGetListingsFromCache_EmptyCache(t *testing.T) {
	// Clear the shared data
	mu.Lock()
	sheetData = nil
	mu.Unlock()

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/listings", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.JSONEq(t, `{"error":"No listings available yet"}`, resp.Body.String())
}

func TestGetSheetDataFromCache_Success(t *testing.T) {
	mu.Lock()
	sheetData = testSheetData
	mu.Unlock()

	data, err := getSheetDataFromCache()
	assert.NoError(t, err)
	assert.Equal(t, testSheetData, data)
}

func TestGetSheetDataFromCache_EmptyCache(t *testing.T) {
	mu.Lock()
	sheetData = nil
	mu.Unlock()

	data, err := getSheetDataFromCache()
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Equal(t, "No listings available yet", err.Error())
}

func TestGetDevListingsFromCache_Success(t *testing.T) {
	mu.Lock()
	devSheetData = testDevSheetData
	mu.Unlock()

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/dev-listings", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header().Get("Content-Type"))
	assert.JSONEq(t, string(testDevSheetData), resp.Body.String())
}

func TestGetDevSheetDataFromCache_EmptyCache(t *testing.T) {
	mu.Lock()
	devSheetData = nil
	mu.Unlock()

	data, err := getDevSheetDataFromCache()
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Equal(t, "No dev listings available yet", err.Error())
}
