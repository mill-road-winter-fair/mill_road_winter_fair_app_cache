package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testSheetData = []byte(`[{"name": "Test Listing"}]`)

// Setup function for tests
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/listings", ListingsEndpoint)
	return router
}

// Test listings endpoint when key is valid and the data is available
func TestGetListingsFromCache_ValidKey_Success(t *testing.T) {
	// Set up test data
	mu.Lock()
	sheetData = testSheetData
	mu.Unlock()

	router := setupRouter()

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/listings", nil)
	req.Header.Add("X-API-Key", "fakeApiKeyForTesting")
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header().Get("Content-Type"))
	assert.JSONEq(t, string(testSheetData), resp.Body.String())
}

// Test listings endpoint when key is invalid and the data is available
func TestGetListingsFromCache_InvalidKey_Success(t *testing.T) {
	// Set up test data
	mu.Lock()
	sheetData = testSheetData
	mu.Unlock()

	router := setupRouter()

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/listings", nil)
	req.Header.Add("X-API-Key", "invalidApiKeyForTesting")
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header().Get("Content-Type"))
	assert.JSONEq(t, string(testSheetData), resp.Body.String())
}

// Test listings endpoint when key is missing and the data is available
func TestGetListingsFromCache_MissingKey_Success(t *testing.T) {
	// Set up test data
	mu.Lock()
	sheetData = testSheetData
	mu.Unlock()

	router := setupRouter()

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/listings", nil)
	req.Header.Add("X-API-Key", "")
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json; charset=UTF-8", resp.Header().Get("Content-Type"))
	assert.JSONEq(t, string(testSheetData), resp.Body.String())
}

// Test listings endpoint when data is not available
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
