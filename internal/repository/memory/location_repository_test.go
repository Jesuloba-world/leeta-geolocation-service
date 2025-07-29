package memory_test

import (
	"fmt"
	"testing"

	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/internal/repository/memory"
)

func TestSave(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()

	// Test saving a new location
	location := &domain.Location{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	err := repo.Save(location)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test saving a duplicate location
	err = repo.Save(location)
	if err == nil {
		t.Error("Expected error for duplicate location, got nil")
	}

	// Test saving nil location
	err = repo.Save(nil)
	if err == nil {
		t.Error("Expected error for nil location, got nil")
	}
}

func TestFindByName(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()

	// Add test location
	location := &domain.Location{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	repo.Save(location)

	// Test finding existing location
	found, err := repo.FindByName("Test Location")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if found.Name != "Test Location" {
		t.Errorf("Expected location name 'Test Location', got '%s'", found.Name)
	}

	// Test finding non-existent location
	_, err = repo.FindByName("Non-existent")
	if err == nil {
		t.Error("Expected error for non-existent location, got nil")
	}

	// Test with empty name
	_, err = repo.FindByName("")
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestFindAll(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()

	// Test with empty repository
	locations, err := repo.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(locations) != 0 {
		t.Errorf("Expected 0 locations, got %d", len(locations))
	}

	// Add test locations
	location1 := &domain.Location{Name: "Location1", Latitude: 40.7128, Longitude: -74.0060}
	location2 := &domain.Location{Name: "Location2", Latitude: 34.0522, Longitude: -118.2437}

	repo.Save(location1)
	repo.Save(location2)

	// Test with populated repository
	locations, err = repo.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(locations))
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()

	// Add test location
	location := &domain.Location{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	repo.Save(location)

	// Test deleting existing location
	err := repo.Delete("Test Location")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify location was deleted
	_, err = repo.FindByName("Test Location")
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}

	// Test deleting non-existent location
	err = repo.Delete("Non-existent")
	if err == nil {
		t.Error("Expected error for non-existent location, got nil")
	}

	// Test with empty name
	err = repo.Delete("")
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			location := &domain.Location{
				Name:      fmt.Sprintf("Location%d", id),
				Latitude:  float64(40 + id),
				Longitude: float64(-74 - id),
			}
			repo.Save(location)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all locations were saved
	locations, err := repo.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(locations) != 10 {
		t.Errorf("Expected 10 locations, got %d", len(locations))
	}
}

func TestRepositoryState(t *testing.T) {
	t.Parallel()
	repo := memory.NewInMemoryLocationRepository()

	// Test initial state
	locations, err := repo.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(locations) != 0 {
		t.Errorf("Expected empty repository, got %d locations", len(locations))
	}

	// Add location
	location := &domain.Location{
		Name:      "Test Location",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	repo.Save(location)

	// Verify state after addition
	locations, err = repo.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(locations))
	}

	// Delete location
	repo.Delete("Test Location")

	// Verify state after deletion
	locations, err = repo.FindAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(locations) != 0 {
		t.Errorf("Expected empty repository after deletion, got %d locations", len(locations))
	}
}