package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"

	"github.com/jesuloba-world/leeta-task/internal/dto"
	"github.com/jesuloba-world/leeta-task/internal/repository/memory"
	"github.com/jesuloba-world/leeta-task/internal/service"
)

func setupTestAPI(t *testing.T) (humatest.TestAPI, *LocationHandler) {
	repo := memory.NewInMemoryLocationRepository()
	locationService := service.NewLocationService(repo)
	locationHandler := NewLocationHandler(locationService)

	_, api := humatest.New(t, huma.DefaultConfig("Test API", "1.0.0"))
	locationHandler.RegisterRoutes(api)

	return api, locationHandler
}

func TestCreateLocation(t *testing.T) {
	api, _ := setupTestAPI(t)

	locationReq := dto.LocationRequest{
		Name:      "New York",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	resp := api.Post("/locations", locationReq)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["name"] != "New York" {
		t.Errorf("Expected name 'New York', got %v", response["name"])
	}
	if response["latitude"] != 40.7128 {
		t.Errorf("Expected latitude 40.7128, got %v", response["latitude"])
	}
	if response["longitude"] != -74.0060 {
		t.Errorf("Expected longitude -74.0060, got %v", response["longitude"])
	}
}

func TestCreateLocationDuplicate(t *testing.T) {
	api, _ := setupTestAPI(t)

	locationReq := dto.LocationRequest{
		Name:      "New York",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	// Create first location
	resp1 := api.Post("/locations", locationReq)
	if resp1.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp1.Code)
	}

	// Try to create duplicate
	resp2 := api.Post("/locations", locationReq)
	if resp2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, resp2.Code)
	}
}

func TestGetAllLocations(t *testing.T) {
	api, _ := setupTestAPI(t)

	locationReq1 := dto.LocationRequest{
		Name:      "New York",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	locationReq2 := dto.LocationRequest{
		Name:      "Los Angeles",
		Latitude:  34.0522,
		Longitude: -118.2437,
	}

	// Create locations
	api.Post("/locations", locationReq1)
	api.Post("/locations", locationReq2)

	// Get all locations
	resp := api.Get("/locations")
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	locations := response["locations"].([]interface{})
	if len(locations) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(locations))
	}
}

func TestDeleteLocation(t *testing.T) {
	api, _ := setupTestAPI(t)

	locationReq := dto.LocationRequest{
		Name:      "To Delete",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	// Create location
	resp1 := api.Post("/locations", locationReq)
	if resp1.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp1.Code)
	}

	// Delete location by name
	resp2 := api.Delete("/locations/To Delete")
	if resp2.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, resp2.Code)
	}

	// Verify deletion
	resp3 := api.Get("/locations")
	if resp3.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp3.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(resp3.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	locations := response["locations"].([]interface{})
	if len(locations) != 0 {
		t.Errorf("Expected 0 locations, got %d", len(locations))
	}
}

func TestDeleteLocationNotFound(t *testing.T) {
	api, _ := setupTestAPI(t)

	// Use a non-existent name
	resp := api.Delete("/locations/non-existent-location")
	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.Code)
	}
}

func TestFindNearest(t *testing.T) {
	api, _ := setupTestAPI(t)

	locationReq1 := dto.LocationRequest{
		Name:      "New York",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	locationReq2 := dto.LocationRequest{
		Name:      "Los Angeles",
		Latitude:  34.0522,
		Longitude: -118.2437,
	}

	// Create locations
	api.Post("/locations", locationReq1)
	api.Post("/locations", locationReq2)

	// Find nearest to a point closer to New York
	resp := api.Get("/nearest?lat=40.7589&lng=-73.9851")
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	location := response["location"].(map[string]interface{})
	if location["name"] != "New York" {
		t.Errorf("Expected location name 'New York', got %v", location["name"])
	}
}

func TestFindNearestMissingParams(t *testing.T) {
	api, _ := setupTestAPI(t)

	resp := api.Get("/nearest?lat=40.7589")
	if resp.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, resp.Code)
	}

	resp = api.Get("/nearest?lng=-73.9851")
	if resp.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, resp.Code)
	}
}

func TestFindNearestNoLocations(t *testing.T) {
	api, _ := setupTestAPI(t)

	resp := api.Get("/nearest?lat=40.7589&lng=-73.9851")
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.Code)
	}
}

func TestCreateLocationInvalidData(t *testing.T) {
	api, _ := setupTestAPI(t)

	tests := []struct {
		name     string
		request  dto.LocationRequest
		expected int
	}{
		{
			name: "empty name",
			request: dto.LocationRequest{
				Name:      "",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expected: 400,
		},
		{
			name: "invalid_latitude",
			request: dto.LocationRequest{
				Name:      "Invalid Lat",
				Latitude:  91.0,
				Longitude: -74.0060,
			},
			expected: 400,
		},
		{
			name: "invalid_longitude",
			request: dto.LocationRequest{
				Name:      "Invalid Lng",
				Latitude:  40.7128,
				Longitude: -181.0,
			},
			expected: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := api.Post("/locations", tt.request)
			if resp.Code != tt.expected {
				t.Errorf("Expected status %d, got %d", tt.expected, resp.Code)
			}
		})
	}
}
