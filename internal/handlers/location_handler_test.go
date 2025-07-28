package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "New York", response["name"])
	assert.Equal(t, 40.7128, response["latitude"])
	assert.Equal(t, -74.0060, response["longitude"])
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
	assert.Equal(t, http.StatusCreated, resp1.Code)

	// Try to create duplicate
	resp2 := api.Post("/locations", locationReq)
	assert.Equal(t, http.StatusConflict, resp2.Code)
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
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)
	locations := response["locations"].([]interface{})
	assert.Len(t, locations, 2)
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
	assert.Equal(t, http.StatusCreated, resp1.Code)

	// Delete location
	resp2 := api.Delete("/locations/To Delete")
	assert.Equal(t, http.StatusNoContent, resp2.Code)

	// Verify deletion
	resp3 := api.Get("/locations")
	assert.Equal(t, http.StatusOK, resp3.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp3.Body.Bytes(), &response)
	require.NoError(t, err)
	locations := response["locations"].([]interface{})
	assert.Len(t, locations, 0)
}

func TestDeleteLocationNotFound(t *testing.T) {
	api, _ := setupTestAPI(t)

	resp := api.Delete("/locations/NonExistent")
	assert.Equal(t, http.StatusNotFound, resp.Code)
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
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)
	location := response["location"].(map[string]interface{})
	assert.Equal(t, "New York", location["name"])
}

func TestFindNearestMissingParams(t *testing.T) {
	api, _ := setupTestAPI(t)

	resp := api.Get("/nearest?lat=40.7589")
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)

	resp = api.Get("/nearest?lng=-73.9851")
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestFindNearestNoLocations(t *testing.T) {
	api, _ := setupTestAPI(t)

	resp := api.Get("/nearest?lat=40.7589&lng=-73.9851")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
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
			assert.Equal(t, tt.expected, resp.Code)
		})
	}
}