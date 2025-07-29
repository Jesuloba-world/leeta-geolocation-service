package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/jesuloba-world/leeta-task/internal/dto"
	"github.com/jesuloba-world/leeta-task/internal/handlers"
	"github.com/jesuloba-world/leeta-task/internal/repository/memory"
	"github.com/jesuloba-world/leeta-task/internal/service"
)

func setupTestServer() http.Handler {
	// Initialize repository
	repo := memory.NewInMemoryLocationRepository()

	// Initialize service
	locationService := service.NewLocationService(repo)

	// Initialize handlers
	locationHandler := handlers.NewLocationHandler(locationService)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create Huma API configuration
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humago.New(mux, config)

	// Register routes with Huma
	locationHandler.RegisterRoutes(api)

	return mux
}

func TestCreateLocation(t *testing.T) {
	t.Parallel()
	server := setupTestServer()

	// Test valid location creation
	locationReq := dto.LocationRequest{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	locationJSON, _ := json.Marshal(locationReq)
	req := httptest.NewRequest("POST", "/locations", bytes.NewBuffer(locationJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	// Test duplicate location
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for duplicate location, got %d", http.StatusBadRequest, rec.Code)
	}

	// Test invalid location (out of range latitude)
	invalidLocation := dto.LocationRequest{
		Name:      "Invalid Location",
		Latitude:  100.0, // Invalid latitude
		Longitude: -74.0060,
	}

	invalidJSON, _ := json.Marshal(invalidLocation)
	req = httptest.NewRequest("POST", "/locations", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid location, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetAllLocations(t *testing.T) {
	t.Parallel()
	server := setupTestServer()

	// Create a test location first
	locationReq := dto.LocationRequest{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	locationJSON, _ := json.Marshal(locationReq)
	req := httptest.NewRequest("POST", "/locations", bytes.NewBuffer(locationJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	// Now test GET /locations
	req = httptest.NewRequest("GET", "/locations", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var response dto.LocationListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response.Locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(response.Locations))
	}

	if response.Locations[0].Name != "Test Location" {
		t.Errorf("Expected location name 'Test Location', got '%s'", response.Locations[0].Name)
	}
}

func TestFindNearest(t *testing.T) {
	t.Parallel()
	server := setupTestServer()

	// Create multiple test locations
	locationReqs := []dto.LocationRequest{
		{
			Name:      "New York",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		{
			Name:      "Los Angeles",
			Latitude:  34.0522,
			Longitude: -118.2437,
		},
		{
			Name:      "Chicago",
			Latitude:  41.8781,
			Longitude: -87.6298,
		},
	}

	for _, loc := range locationReqs {
		locationJSON, _ := json.Marshal(loc)
		req := httptest.NewRequest("POST", "/locations", bytes.NewBuffer(locationJSON))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
	}

	// Test finding nearest to a point near Chicago
	req := httptest.NewRequest("GET", "/nearest?lat=42.0&lng=-88.0", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var response dto.NearestLocationResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Location.Name != "Chicago" {
		t.Errorf("Expected nearest location to be 'Chicago', got '%s'", response.Location.Name)
	}
}

func TestDeleteLocation(t *testing.T) {
	t.Parallel()
	server := setupTestServer()

	// Create a test location first
	locationReq := dto.LocationRequest{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	locationJSON, _ := json.Marshal(locationReq)
	req := httptest.NewRequest("POST", "/locations", bytes.NewBuffer(locationJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	// Now test DELETE /locations/{name}
	req = httptest.NewRequest("DELETE", "/locations/Test%20Location", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Verify location is deleted
	req = httptest.NewRequest("GET", "/locations", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var response dto.LocationListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response.Locations) != 0 {
		t.Errorf("Expected 0 locations after deletion, got %d", len(response.Locations))
	}
}

func TestAPIErrorHandling(t *testing.T) {
	t.Parallel()
	server := setupTestServer()

	// Test invalid JSON
	req := httptest.NewRequest("POST", "/locations", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid JSON, got %d", http.StatusBadRequest, rec.Code)
	}

	// Test missing Content-Type
	locationReq := dto.LocationRequest{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	locationJSON, _ := json.Marshal(locationReq)
	req = httptest.NewRequest("POST", "/locations", bytes.NewBuffer(locationJSON))
	// Don't set Content-Type header

	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d for missing Content-Type, got %d", http.StatusCreated, rec.Code)
	}

	// Test invalid query parameters for nearest endpoint
	req = httptest.NewRequest("GET", "/nearest?lat=invalid&lng=-88.0", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d for invalid lat parameter, got %d", http.StatusUnprocessableEntity, rec.Code)
	}

	// Test missing query parameters for nearest endpoint
	req = httptest.NewRequest("GET", "/nearest", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d for missing parameters, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()
	server := setupTestServer()

	// Test unsupported method on /locations
	req := httptest.NewRequest("PUT", "/locations", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for unsupported method, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	// Test unsupported method on /nearest
	req = httptest.NewRequest("POST", "/nearest", nil)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for unsupported method, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}