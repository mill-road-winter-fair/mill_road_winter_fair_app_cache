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
	router.GET("/listings", ListingsEndpoint)
	router.GET("/dev-listings", DevListingsEndpoint)
	router.GET("/notifications", NotificationsEndpoint)
	router.GET("/notifications-dev", DevNotificationsEndpoint)
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

// Test both public routes with distinct caches and every supported credential source.
func TestNotificationsEndpoints(t *testing.T) {
	previousKey := ourApiKey
	ourApiKey = "test-api-key"
	mu.Lock()
	previousProd, previousDev := notificationsSheetData, devNotificationsSheetData
	mu.Unlock()
	t.Cleanup(func() {
		ourApiKey = previousKey
		mu.Lock()
		notificationsSheetData, devNotificationsSheetData = previousProd, previousDev
		mu.Unlock()
	})

	for _, empty := range []bool{false, true} {
		mu.Lock()
		notificationsSheetData = []byte(`{"values":[["Production notification"]]}`)
		devNotificationsSheetData = []byte(`{"values":[["Development notification"]]}`)
		if empty {
			notificationsSheetData, devNotificationsSheetData = nil, nil
		}
		mu.Unlock()
		for _, route := range []struct {
			path string
			body string
		}{
			{"/notifications", `{"values":[["Production notification"]]}`},
			{"/notifications-dev", `{"values":[["Development notification"]]}`},
		} {
			for _, credential := range []struct {
				name       string
				header     string
				query      string
				authorised bool
			}{
				{"missing", "", "", false},
				{"invalid header", "wrong", "", false},
				{"invalid query", "", "?key=wrong", false},
				{"valid header", "test-api-key", "", true},
				{"valid query", "", "?key=test-api-key", true},
				{"invalid header overrides query", "wrong", "?key=test-api-key", false},
				{"valid header overrides query", "test-api-key", "?key=wrong", true},
			} {
				t.Run(route.path+"/"+credential.name+map[bool]string{true: "/empty", false: "/populated"}[empty], func(t *testing.T) {
					req := httptest.NewRequest(http.MethodGet, route.path+credential.query, nil)
					if credential.header != "" {
						req.Header.Set("X-API-Key", credential.header)
					}
					resp := httptest.NewRecorder()
					setupRouter().ServeHTTP(resp, req)
					switch {
					case !credential.authorised:
						assert.Equal(t, http.StatusUnauthorized, resp.Code)
						assert.JSONEq(t, `{"error":"missing or invalid API key"}`, resp.Body.String())
					case empty:
						assert.Equal(t, http.StatusInternalServerError, resp.Code)
						assert.JSONEq(t, `{"error":"No notifications available yet"}`, resp.Body.String())
					default:
						assert.Equal(t, http.StatusOK, resp.Code)
						assert.JSONEq(t, route.body, resp.Body.String())
					}
				})
			}
		}
	}
	ourApiKey = ""
	for _, path := range []string{"/notifications", "/notifications-dev"} {
		resp := httptest.NewRecorder()
		setupRouter().ServeHTTP(resp, httptest.NewRequest(http.MethodGet, path+"?key=test-api-key", nil))
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}

func TestListingsRemainAccessibleWithInvalidKey(t *testing.T) {
	mu.Lock()
	sheetData, devSheetData = testSheetData, testDevSheetData
	mu.Unlock()
	for _, path := range []string{"/listings", "/dev-listings"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("X-API-Key", "invalid")
		resp := httptest.NewRecorder()
		setupRouter().ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}
