package service_test

import (
	"testing"

	"github.com/jesuloba-world/leeta-task/internal/repository/memory"
	"github.com/jesuloba-world/leeta-task/internal/service"
)

func TestCreateLocation(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()
	svc := service.NewLocationService(repo)

	// Test valid location creation
	location, err := svc.CreateLocation("Test Location", 40.7128, -74.0060)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if location == nil {
		t.Error("Expected location to be created, got nil")
	}

	// Test duplicate location
	_, err = svc.CreateLocation("Test Location", 40.7128, -74.0060)
	if err == nil {
		t.Error("Expected error for duplicate location, got nil")
	}

	// Test invalid location (empty name)
	_, err = svc.CreateLocation("", 40.7128, -74.0060)
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestGetAllLocations(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()
	svc := service.NewLocationService(repo)

	// Create test locations
	_, err := svc.CreateLocation("Location1", 40.7128, -74.0060)
	if err != nil {
		t.Errorf("Expected no error creating location 1, got %v", err)
	}
	_, err = svc.CreateLocation("Location2", 34.0522, -118.2437)
	if err != nil {
		t.Errorf("Expected no error creating location 2, got %v", err)
	}

	locations, err := svc.GetAllLocations()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(locations))
	}
}

func TestGetLocationByName(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()
	svc := service.NewLocationService(repo)

	// Create test location
	_, err := svc.CreateLocation("Test Location", 40.7128, -74.0060)
	if err != nil {
		t.Errorf("Expected no error creating location, got %v", err)
	}

	// Test existing location
	found, err := svc.GetLocation("Test Location")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if found.Name != "Test Location" {
		t.Errorf("Expected location name 'Test Location', got '%s'", found.Name)
	}

	// Test non-existent location
	_, err = svc.GetLocation("Non-existent")
	if err == nil {
		t.Error("Expected error for non-existent location, got nil")
	}
}

func TestDeleteLocation(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()
	svc := service.NewLocationService(repo)

	// Create a test location first
	_, err := svc.CreateLocation("Test Location", 40.7128, -74.0060)
	if err != nil {
		t.Errorf("Expected no error creating location, got %v", err)
	}

	// Test deleting existing location
	err = svc.DeleteLocation("Test Location")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify location was deleted
	_, err = svc.GetLocation("Test Location")
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}

	// Test deleting non-existent location
	err = svc.DeleteLocation("Non-existent")
	if err == nil {
		t.Error("Expected error for non-existent location, got nil")
	}
}

func TestFindNearest(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()
	svc := service.NewLocationService(repo)

	// Create test locations
	testLocations := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"New York", 40.7128, -74.0060},
		{"Los Angeles", 34.0522, -118.2437},
		{"Chicago", 41.8781, -87.6298},
	}

	for _, loc := range testLocations {
		_, err := svc.CreateLocation(loc.name, loc.lat, loc.lng)
		if err != nil {
			t.Errorf("Expected no error creating location %s, got %v", loc.name, err)
		}
	}

	// Test finding nearest to a point near Chicago
	nearest, distance, err := svc.FindNearest(42.0, -88.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if nearest.Name != "Chicago" {
		t.Errorf("Expected nearest location to be 'Chicago', got '%s'", nearest.Name)
	}

	if distance <= 0 {
		t.Errorf("Expected positive distance, got %f", distance)
	}

	// Test with empty repository
	emptyRepo := memory.NewInMemoryLocationRepository()
	emptySvc := service.NewLocationService(emptyRepo)

	_, _, err = emptySvc.FindNearest(42.0, -88.0)
	if err == nil {
		t.Error("Expected error with empty repository, got nil")
	}
}

func TestCreateLocationValidation(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()
	svc := service.NewLocationService(repo)

	tests := []struct {
		name      string
		location  string
		latitude  float64
		longitude float64
		wantErr   bool
	}{
		{
			name:      "Valid location",
			location:  "Valid Location",
			latitude:  40.7128,
			longitude: -74.0060,
			wantErr:   false,
		},
		{
			name:      "Empty name",
			location:  "",
			latitude:  40.7128,
			longitude: -74.0060,
			wantErr:   true,
		},
		{
			name:      "Latitude too high",
			location:  "Invalid Lat High",
			latitude:  91.0,
			longitude: -74.0060,
			wantErr:   true,
		},
		{
			name:      "Latitude too low",
			location:  "Invalid Lat Low",
			latitude:  -91.0,
			longitude: -74.0060,
			wantErr:   true,
		},
		{
			name:      "Longitude too high",
			location:  "Invalid Lng High",
			latitude:  40.7128,
			longitude: 181.0,
			wantErr:   true,
		},
		{
			name:      "Longitude too low",
			location:  "Invalid Lng Low",
			latitude:  40.7128,
			longitude: -181.0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateLocation(tt.location, tt.latitude, tt.longitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}