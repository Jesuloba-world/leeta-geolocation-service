package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesuloba-world/leeta-task/internal/domain"
)

// MockLocationService implements the LocationService interface for testing
type MockLocationService struct {
	locations map[string]*domain.Location
	createError error
	getAllError error
	deleteError error
	findNearestError error
}

func NewMockLocationService() *MockLocationService {
	return &MockLocationService{
		locations: make(map[string]*domain.Location),
	}
}

func (m *MockLocationService) CreateLocation(name string, latitude, longitude float64) (*domain.Location, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	if _, exists := m.locations[name]; exists {
		return nil, errors.New("location already exists")
	}
	location, err := domain.NewLocation(name, latitude, longitude)
	if err != nil {
		return nil, err
	}
	m.locations[name] = location
	return location, nil
}

func (m *MockLocationService) GetLocation(name string) (*domain.Location, error) {
	location, exists := m.locations[name]
	if !exists {
		return nil, errors.New("location not found")
	}
	return location, nil
}

func (m *MockLocationService) GetAllLocations() ([]*domain.Location, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}
	locations := make([]*domain.Location, 0, len(m.locations))
	for _, location := range m.locations {
		locations = append(locations, location)
	}
	return locations, nil
}

func (m *MockLocationService) DeleteLocation(name string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	if _, exists := m.locations[name]; !exists {
		return errors.New("location not found")
	}
	delete(m.locations, name)
	return nil
}

func (m *MockLocationService) FindNearest(lat, lng float64) (*domain.Location, float64, error) {
	if m.findNearestError != nil {
		return nil, 0, m.findNearestError
	}
	if len(m.locations) == 0 {
		return nil, 0, errors.New("no locations available")
	}
	// Return the first location with a mock distance
	for _, location := range m.locations {
		return location, 10.5, nil
	}
	return nil, 0, errors.New("no locations available")
}

func TestCreateLocation(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		setupMock      func(*MockLocationService)
	}{
		{
			name:           "Valid location creation",
			method:         "POST",
			body:           `{"name":"Test Location","latitude":40.7128,"longitude":-74.0060}`,
			expectedStatus: http.StatusCreated,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Invalid method",
			method:         "GET",
			body:           `{"name":"Test Location","latitude":40.7128,"longitude":-74.0060}`,
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Invalid JSON body",
			method:         "POST",
			body:           `{"invalid":"json"}`,
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Duplicate location",
			method:         "POST",
			body:           `{"name":"Existing Location","latitude":40.7128,"longitude":-74.0060}`,
			expectedStatus: http.StatusConflict,
			setupMock: func(m *MockLocationService) {
				m.createError = errors.New("location already exists")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockLocationService()
			tt.setupMock(mockService)
			handler := NewLocationHandler(mockService)

			req := httptest.NewRequest(tt.method, "/locations", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			handler.CreateLocation(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetAllLocations(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		setupMock      func(*MockLocationService)
	}{
		{
			name:           "Valid get all locations",
			method:         "GET",
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockLocationService) {
				location, _ := domain.NewLocation("Test", 40.7128, -74.0060)
				m.locations["Test"] = location
			},
		},
		{
			name:           "Invalid method",
			method:         "POST",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Service error",
			method:         "GET",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *MockLocationService) {
				m.getAllError = errors.New("service error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockLocationService()
			tt.setupMock(mockService)
			handler := NewLocationHandler(mockService)

			req := httptest.NewRequest(tt.method, "/locations", nil)
			w := httptest.NewRecorder()

			handler.GetAllLocations(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteLocation(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMock      func(*MockLocationService)
	}{
		{
			name:           "Valid location deletion",
			method:         "DELETE",
			path:           "/locations/test",
			expectedStatus: http.StatusNoContent,
			setupMock: func(m *MockLocationService) {
				location, _ := domain.NewLocation("test", 40.7128, -74.0060)
				m.locations["test"] = location
			},
		},
		{
			name:           "Invalid method",
			method:         "GET",
			path:           "/locations/test",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Invalid path",
			method:         "DELETE",
			path:           "/locations",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Location not found",
			method:         "DELETE",
			path:           "/locations/nonexistent",
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockLocationService) {
				m.deleteError = errors.New("location not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockLocationService()
			tt.setupMock(mockService)
			handler := NewLocationHandler(mockService)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.DeleteLocation(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestFindNearest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMock      func(*MockLocationService)
	}{
		{
			name:           "Valid find nearest",
			method:         "GET",
			path:           "/nearest?lat=40.7128&lng=-74.0060",
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockLocationService) {
				location, _ := domain.NewLocation("test", 40.7128, -74.0060)
				m.locations["test"] = location
			},
		},
		{
			name:           "Invalid method",
			method:         "POST",
			path:           "/nearest?lat=40.7128&lng=-74.0060",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Missing latitude",
			method:         "GET",
			path:           "/nearest?lng=-74.0060",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Missing longitude",
			method:         "GET",
			path:           "/nearest?lat=40.7128",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Invalid latitude",
			method:         "GET",
			path:           "/nearest?lat=invalid&lng=-74.0060",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "Invalid longitude",
			method:         "GET",
			path:           "/nearest?lat=40.7128&lng=invalid",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockLocationService) {},
		},
		{
			name:           "No locations available",
			method:         "GET",
			path:           "/nearest?lat=40.7128&lng=-74.0060",
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockLocationService) {
				m.findNearestError = errors.New("no locations available")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockLocationService()
			tt.setupMock(mockService)
			handler := NewLocationHandler(mockService)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.FindNearest(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestLocationHandlerIntegration(t *testing.T) {
	mockService := NewMockLocationService()
	handler := NewLocationHandler(mockService)

	// Test creating a location
	createBody := `{"name":"NYC","latitude":40.7128,"longitude":-74.0060}`
	createReq := httptest.NewRequest("POST", "/locations", bytes.NewBufferString(createBody))
	createW := httptest.NewRecorder()
	handler.CreateLocation(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Errorf("Expected status %d for create, got %d", http.StatusCreated, createW.Code)
	}

	// Test getting all locations
	getAllReq := httptest.NewRequest("GET", "/locations", nil)
	getAllW := httptest.NewRecorder()
	handler.GetAllLocations(getAllW, getAllReq)

	if getAllW.Code != http.StatusOK {
		t.Errorf("Expected status %d for get all, got %d", http.StatusOK, getAllW.Code)
	}

	var locations []*domain.Location
	if err := json.NewDecoder(getAllW.Body).Decode(&locations); err != nil {
		t.Errorf("Failed to decode locations: %v", err)
	}

	if len(locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(locations))
	}

	// Test finding nearest location
	nearestReq := httptest.NewRequest("GET", "/nearest?lat=40.7128&lng=-74.0060", nil)
	nearestW := httptest.NewRecorder()
	handler.FindNearest(nearestW, nearestReq)

	if nearestW.Code != http.StatusOK {
		t.Errorf("Expected status %d for find nearest, got %d", http.StatusOK, nearestW.Code)
	}

	// Test deleting the location
	deleteReq := httptest.NewRequest("DELETE", "/locations/NYC", nil)
	deleteW := httptest.NewRecorder()
	handler.DeleteLocation(deleteW, deleteReq)

	if deleteW.Code != http.StatusNoContent {
		t.Errorf("Expected status %d for delete, got %d", http.StatusNoContent, deleteW.Code)
	}
}